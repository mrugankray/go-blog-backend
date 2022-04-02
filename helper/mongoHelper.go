package helper

import (
	"context"
	"errors"
	"fmt"

	"com.mrugankray.blog/conf"
	"com.mrugankray.blog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// collection instances
var postCollection *mongo.Collection
var userCollection *mongo.Collection

func init() {
	postCollection, userCollection = conf.ConnectToDB()
}

// POSTS

func InsertPostHelper(post *model.Post) (*mongo.InsertOneResult, error) {
	insertedPost, err := postCollection.InsertOne(context.Background(), *post)

	if err != nil {
		return nil, err
	}

	return insertedPost, nil
}

func UpdatePostHelper(postId string, post *model.Post) (*mongo.UpdateResult, error) {
	// query
	updateQuery := bson.M{"$set": bson.M{"date": post.Date, "title": post.Title, "body": post.Body, "userId": post.UserId}}

	// convert string to object id
	id, objIdErr := primitive.ObjectIDFromHex(postId)
	if objIdErr != nil {
		return nil, errors.New("invalid object id")
	}
	filter := bson.M{"_id": id}

	// update
	updatedPost, err := postCollection.UpdateOne(context.Background(), filter, updateQuery)

	if err != nil {
		return nil, errors.New("something went wrong while updating post")
	}

	return updatedPost, nil
}

func GetPostsHelper() ([]*model.Post, error) {
	cursor, err := postCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		return nil, err
	}

	var posts []*model.Post

	for cursor.Next(context.Background()) {
		var post model.Post
		err := cursor.Decode(&post)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func DeletePostHelper(postId string) (*mongo.DeleteResult, error) {
	// convert string to object id
	id, objIdErr := primitive.ObjectIDFromHex(postId)
	if objIdErr != nil {
		return nil, errors.New("invalid object id")
	}
	filter := bson.M{"_id": id}

	// delete
	deleteResult, err := postCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, errors.New("something went wrong while deleting post")
	}

	return deleteResult, nil
}

// USERS

func InsertUserHelper(user *model.User) (*mongo.InsertOneResult, error) {
	insertedUser, err := userCollection.InsertOne(context.Background(), *user)

	if err != nil {
		return nil, err
	}

	return insertedUser, nil
}

func UpdateUserHelper(userId string, user *model.User) (*mongo.UpdateResult, error) {
	fieldsToBeUpdated := make(map[string]string)

	// check what fields are present
	if user.Email != "" {
		fieldsToBeUpdated["email"] = user.Email
	}
	if user.Password != "" {
		fieldsToBeUpdated["password"] = user.Password
	}
	if user.FullName != "" {
		fieldsToBeUpdated["fullName"] = user.FullName
	}

	// query
	updateQuery := bson.M{"$set": fieldsToBeUpdated}

	// convert string to object id
	id, objIdErr := primitive.ObjectIDFromHex(userId)
	if objIdErr != nil {
		return nil, errors.New("invalid object id")
	}
	filter := bson.M{"_id": id}

	// update
	updatedUser, err := userCollection.UpdateOne(context.Background(), filter, updateQuery)

	if err != nil {
		return nil, errors.New("something went wrong while updating user")
	}

	return updatedUser, nil
}

func GetUsersHelper() ([]*model.User, error) {
	cursor, err := userCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		return nil, err
	}

	var users []*model.User

	for cursor.Next(context.Background()) {
		var user model.User
		err := cursor.Decode(&user)

		if err != nil {
			return nil, err
		}

		// remove password from user
		user.Password = ""

		users = append(users, &user)
	}

	return users, nil
}

func GetUserByEmailHelper(email string, user *model.User) (*model.User, error) {
	filter := bson.M{"email": email}

	// Find user
	userResult := userCollection.FindOne(context.Background(), filter)

	if userResult.Err() != nil {
		return nil, fmt.Errorf("no user with %s email found", email)
	}

	// decode
	if decodeError := userResult.Decode(user); decodeError != nil {
		fmt.Println(decodeError)
		return nil, errors.New("something went wrong while get user details")
	}

	return user, nil
}

func GetUserByIdHelper(userId string, user *model.User) (*model.User, error) {
	// convert string to object id
	id, userIdErr := primitive.ObjectIDFromHex(userId)
	if userIdErr != nil {
		return nil, errors.New("invalid object id")
	}
	filter := bson.M{"_id": id}

	// Find user
	userResult := userCollection.FindOne(context.Background(), filter)

	if userResult.Err() != nil {
		return nil, fmt.Errorf("no user found")
	}

	// decode
	if decodeError := userResult.Decode(user); decodeError != nil {
		fmt.Println(decodeError)
		return nil, errors.New("something went wrong while getting user details")
	}

	return user, nil
}

func DeleteUserHelper(userId string) (*mongo.DeleteResult, error) {
	// convert string to object id
	id, objIdErr := primitive.ObjectIDFromHex(userId)
	if objIdErr != nil {
		return nil, errors.New("invalid object id")
	}
	filter := bson.M{"_id": id}

	// delete
	deleteResult, err := userCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, errors.New("something went wrong while deleting user")
	}

	return deleteResult, nil
}
