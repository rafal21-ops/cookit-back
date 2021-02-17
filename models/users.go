package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CollectionUsers = "users"
var ErrCreateUser = errors.New("Failed to create user record")
var ErrFindUser = errors.New("Failed to find user record")

type User struct {
	ID          string `json:"id" bson:"id"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	AvatarURL   string `json:"avatarURL" bson:"avatar_url"`
	Email       string `json:"email" bson:"email"`
	Description string `json:"description" bson:"description"`
}
type Credentials struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
type Login struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func RegisterUser(client *mongo.Client, db string, u *User) error {
	fmt.Println("User:", u, "Database name:", db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionUsers)
	_, err := collection.InsertOne(ctx, &u)
	if err != nil {
		return ErrCreateUser
	}
	return nil
}
func GetPasswordByUsernameOrEmail(client *mongo.Client, db string, username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionUsers)
	result := User{}
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&result); err == nil {
		return result.Password, nil
	}
	if err := collection.FindOne(ctx, bson.M{"email": username}).Decode(&result); err == nil {
		return result.Password, nil
	}
	return "", ErrFindUser
}
