package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"com.mrugankray.blog/constants"
	"com.mrugankray.blog/helper"
	"com.mrugankray.blog/middleware"
	"com.mrugankray.blog/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// import env
	var err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorResponse := make(map[string][]map[string]string)

	// validate request body
	user, validateBodyError := middleware.ValidateCreateUserBodyMiddleware(w, r)
	if validateBodyError != nil {
		return
	}

	// check if there is no other user with the same email
	_, userByEmailError := helper.GetUserByEmailHelper(user.Email, user)
	// if there is a user with same email then the helper func returns empty error
	if userByEmailError == nil {
		w.WriteHeader(http.StatusForbidden)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "Incorrect email/password"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// generate id for user
	user.ID = primitive.NewObjectID()

	// hash password
	password, hashPasswordError := bcrypt.GenerateFromPassword([]byte(user.Password), constants.HashCost)
	user.Password = string(password)
	if hashPasswordError != nil {
		fmt.Println(hashPasswordError)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// insert user into db
	_, insertErr := helper.InsertUserHelper(user)
	if insertErr != nil {
		fmt.Println(insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// remove password in json response
	user.Password = ""

	json.NewEncoder(w).Encode(*user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorResponse := make(map[string][]map[string]string)

	// validate request body
	user, validateBodyError := middleware.ValidateLoginUserBodyMiddleware(w, r)
	if validateBodyError != nil {
		return
	}

	// keep request password in a new var
	passwordInRequestBody := user.Password

	// get user by email
	_, getUserError := helper.GetUserByEmailHelper(user.Email, user)
	if getUserError != nil {
		// fmt.Println(getUserError)
		w.WriteHeader(http.StatusNotFound)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": getUserError.Error()})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// compare password
	isMatchError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordInRequestBody))

	if isMatchError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "Incorrect email/password"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	mySigningKey := []byte(os.Getenv("SECRET"))

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Add(time.Hour*time.Duration(24)).Unix(), 0)),
		Issuer:    user.ID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, signedTokenError := token.SignedString(mySigningKey)

	if signedTokenError != nil {
		fmt.Println(signedTokenError)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "Something went wrong"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	json.NewEncoder(w).Encode(signedToken)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string][]map[string]string)

	// validate token
	userId, validateTokenError := middleware.ValidateTokenInHeader(w, r)
	if validateTokenError != nil {
		return
	}

	var user model.User

	// get user by userId
	_, getUserError := helper.GetUserByIdHelper(*userId, &user)
	if getUserError != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": getUserError.Error()})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// remove password in json response
	user.Password = ""

	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string][]map[string]string)

	users, getUsersError := helper.GetUsersHelper()
	if getUsersError != nil {
		fmt.Println(getUsersError)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// check if users is empty
	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "No users found"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// Get a user
func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user model.User

	// get user
	_, err := helper.GetUserByIdHelper(params["id"], &user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// remove password
	user.Password = ""

	json.NewEncoder(w).Encode(user)
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorResponse := make(map[string][]map[string]string)
	var user model.User

	// validate token
	userId, validateTokenError := middleware.ValidateTokenInHeader(w, r)
	if validateTokenError != nil {
		return
	}

	// decode
	decodeError := json.NewDecoder(r.Body).Decode(&user)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "Invalid request body"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// hash password
	if user.Password != "" {
		password, hashPasswordError := bcrypt.GenerateFromPassword([]byte(user.Password), constants.HashCost)
		user.Password = string(password)
		if hashPasswordError != nil {
			fmt.Println(hashPasswordError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Something went wrong")
			return
		}
	}

	// update user
	updatResult, updatedUserError := helper.UpdateUserHelper(*userId, &user)
	if updatedUserError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(updatedUserError.Error())
		return
	}

	if updatResult.ModifiedCount > 0 {
		// set ID to id from params
		user.ID, _ = primitive.ObjectIDFromHex(*userId)

		// remove password from response
		user.Password = ""

		json.NewEncoder(w).Encode(user)
	} else {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "no such user found"})

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}
}

// delete a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	errorResponse := make(map[string][]map[string]string)

	// validate token
	userId, validateTokenError := middleware.ValidateTokenInHeader(w, r)
	if validateTokenError != nil {
		return
	}

	// delete
	deleteResult, deleteError := helper.DeleteUserHelper(*userId)

	if deleteError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(deleteError.Error())
		return
	}

	if deleteResult.DeletedCount > 0 {
		responseJson := map[string]string{"deletedId": *userId}

		json.NewEncoder(w).Encode(responseJson)
	} else {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "post with this id is not found"})

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}
}
