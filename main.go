package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"com.mrugankray.blog/router"
	"github.com/joho/godotenv"
)

// import env
var err = godotenv.Load(".env")

func main() {
	// checks if any error occured during import of env file
	if err != nil {
		log.Fatal(err)
	}

	r := router.Router()
	fmt.Printf("Listening on port%s \n", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), r))
}
