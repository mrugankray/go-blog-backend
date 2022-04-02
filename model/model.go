package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Date   string             `json:"date"`
	Title  string             `json:"title"`
	Body   string             `json:"body"`
	UserId primitive.ObjectID `json:"userId" bson:"userId"`
}

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Email    string             `json:"email"`
	Password string             `json:"password,omitempty"`
	FullName string             `json:"fullName,omitempty"`
}
