package conf

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// import env
var err = godotenv.Load(".env")

var connectionString = "mongodb+srv://{username}:{password}@go-learn.i9qdc.mongodb.net/{db}?retryWrites=true&w=majority"

func ConnectToDB() (*mongo.Collection, *mongo.Collection) {
	// sudo gedit /etc/resolv.conf
	// 127.0.0.53

	// mongo password url encoded
	dbPassword := url.QueryEscape(os.Getenv("DB_PASSWORD"))

	connectionString = strings.Replace(connectionString, "{username}", os.Getenv("DB_USERNAME"), 1)
	connectionString = strings.Replace(connectionString, "{password}", dbPassword, 1)
	connectionString = strings.Replace(connectionString, "{db}", os.Getenv("DB_NAME"), 1)

	// options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		// return err
	}

	fmt.Println("Connected to mongo")

	// collections
	postCollection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("POST_COLLECTION_NAME"))
	userCollection := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USER_COLLECTION_NAME"))

	return postCollection, userCollection
}
