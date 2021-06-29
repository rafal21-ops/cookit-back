package models

import (
	"context"
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CollectionRecipes = "recipes"
var ErrCreateRecipe = errors.New("Failed to create recipe record")
var ErrFindRecipe = errors.New("Failed to find recipe record")

type KitchenStyle int

const (
	Polish KitchenStyle = iota
	Russian
)

type Tags int

const (
	Easy Tags = iota
	Medium
	Hard
)

type Recipe struct {
	ID               string       `json:"id" bson:"id"`
	Name             string       `json:"name" bson:"name"`
	UserID           string       `json:"userId" bson:"user_id"`
	Kitchen          KitchenStyle `json:"kitchenStyle" bson:"kitchen_style"`
	Tags             Tags         `json:"tags" bson:"tags"`
	ListOfSteps      []string     `json:"listOfSteps" bson:"list_of_steps"`
	ListOfCategories []Category   `json:"listOfCategories" bson:"list_of_categories"`
	Description      string       `json:"description" bson:"description"`
}

func CreateRecipe(client *mongo.Client, db string, r *Recipe) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionRecipes)
	if _, err := collection.InsertOne(ctx, &r); err != nil {
		return ErrCreateRecipe
	}
	return nil
}
func GetAllRecipes(client *mongo.Client, db string) ([]Recipe, error) {
	recipes := []Recipe{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionRecipes)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var recipe Recipe
		if err = cursor.Decode(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
