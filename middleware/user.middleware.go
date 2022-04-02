package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"com.mrugankray.blog/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func init() {
	// import env
	var err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateCreateUserBodyMiddleware(w http.ResponseWriter, r *http.Request) (*model.User, error) {
	var user model.User

	errorResponse := make(map[string][]map[string]string)

	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "error while parsing request body"})
		json.NewEncoder(w).Encode(errorResponse)
		return nil, decodeErr
	}

	// check for empty fields
	if user.Email == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "email can't be empty", "param": "title"})
	}
	if user.Password == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "password can't be empty", "param": "body"})
	}
	if user.FullName == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "full name can't be empty", "param": "date"})
	}

	if len(errorResponse["errors"]) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return nil, errors.New("errors in request body")
	}

	return &user, nil
}

func ValidateLoginUserBodyMiddleware(w http.ResponseWriter, r *http.Request) (*model.User, error) {
	var user model.User

	errorResponse := make(map[string][]map[string]string)

	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "error while parsing request body"})
		json.NewEncoder(w).Encode(errorResponse)
		return nil, decodeErr
	}

	// check for empty fields
	if user.Email == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "email can't be empty", "param": "title"})
	}
	if user.Password == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "password can't be empty", "param": "body"})
	}

	if len(errorResponse["errors"]) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return nil, errors.New("errors in request body")
	}

	return &user, nil
}

func ValidateTokenInHeader(w http.ResponseWriter, r *http.Request) (*string, error) {
	token := r.Header.Get("Authorization")

	errorResponse := make(map[string][]map[string]string)

	// check if token is empty
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "You're not authorized"})
		json.NewEncoder(w).Encode(errorResponse)
		return nil, errors.New("you're not authorized")
	}

	// parse token
	parsedToken, parseTokenError := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if parseTokenError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Something went wrong")
		return nil, errors.New("something went wrong")
	}

	// get user id
	claims := parsedToken.Claims.(*jwt.RegisteredClaims)
	return &claims.Issuer, nil
}
