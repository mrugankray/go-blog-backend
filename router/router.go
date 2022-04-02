package router

import (
	"com.mrugankray.blog/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// posts
	router.HandleFunc("/api/posts", controller.CreatePost).Methods("POST")
	router.HandleFunc("/api/posts", controller.GetPosts).Methods("GET")
	router.HandleFunc("/api/posts/{id}", controller.UpdatePost).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", controller.DeletePost).Methods("DELETE")

	// users
	router.HandleFunc("/api/user", controller.GetUserDetails).Methods("GET")
	router.HandleFunc("/api/user/{id}", controller.GetUserById).Methods("GET")
	router.HandleFunc("/api/user/all", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/user/signup", controller.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", controller.LoginUser).Methods("POST")
	router.HandleFunc("/api/user", controller.UpdateUserDetails).Methods("PUT")
	router.HandleFunc("/api/user", controller.DeleteUser).Methods("DELETE")

	// middleware
	router.Use(mux.CORSMethodMiddleware(router))

	return router
}
