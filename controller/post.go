package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"com.mrugankray.blog/helper"
	"com.mrugankray.blog/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// inserts a single post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// validate token
	userId, validateTokenError := middleware.ValidateTokenInHeader(w, r)
	if validateTokenError != nil {
		return
	}

	// validate body
	post, validateBodyError := middleware.ValidatePostBodyMiddleware(w, r)
	if validateBodyError != nil {
		return
	}

	// generte doc id
	post.ID = primitive.NewObjectID()
	post.UserId, _ = primitive.ObjectIDFromHex(*userId)

	// insert into db
	_, insertErr := helper.InsertPostHelper(post)
	if insertErr != nil {
		fmt.Println(insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	json.NewEncoder(w).Encode(*post)
}

// get all posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorResponse := make(map[string][]map[string]string)

	posts, err := helper.GetPostsHelper()

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// check if posts is empty
	if len(posts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "No posts found"})
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

// delete a post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	errorResponse := make(map[string][]map[string]string)

	// delete
	deleteResult, err := helper.DeletePostHelper(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if deleteResult.DeletedCount > 0 {
		responseJson := map[string]string{"deletedId": params["id"]}

		json.NewEncoder(w).Encode(responseJson)
	} else {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "post with this id is not found"})

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}
}

// update a post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	errorResponse := make(map[string][]map[string]string)

	// validate body
	post, validateBodyError := middleware.ValidatePostBodyMiddleware(w, r)
	if validateBodyError != nil {
		return
	}

	// update
	updatResult, err := helper.UpdatePostHelper(params["id"], post)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if updatResult.ModifiedCount > 0 {
		// set ID to id from params
		post.ID, _ = primitive.ObjectIDFromHex(params["id"])

		json.NewEncoder(w).Encode(*post)
	} else {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "post with this id is not found"})

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}
}
