package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"com.mrugankray.blog/model"
)

func ValidatePostBodyMiddleware(w http.ResponseWriter, r *http.Request) (*model.Post, error) {
	var post model.Post

	errorResponse := make(map[string][]map[string]string)

	decodeErr := json.NewDecoder(r.Body).Decode(&post)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		w.WriteHeader(http.StatusBadRequest)
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "error while parsing request body"})
		json.NewEncoder(w).Encode(errorResponse)
		return nil, decodeErr
	}

	// check for empty fields
	if post.Title == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "title can't be empty", "param": "title"})
	}
	if post.Body == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "body can't be empty", "param": "body"})
	}
	if post.Date == "" {
		errorResponse["errors"] = append(errorResponse["errors"], map[string]string{"msg": "date can't be empty", "param": "date"})
	}

	// errorResponseMarshaled, _ := json.Marshal(errorResponse)
	// fmt.Println(string(errorResponseMarshaled))

	if len(errorResponse["errors"]) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return nil, errors.New("errors in request body")
	}

	return &post, nil
}
