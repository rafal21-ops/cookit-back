package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Name             string       `json:"name" bson:"name" validate:"nonnil,nonzero"`
	UserID           string       `json:"userId" bson:"user_id" validate:"nonnil,nonzero"`
	Kitchen          KitchenStyle `json:"kitchenStyle" bson:"kitchen_style" validate:"nonnil,nonzero"`
	Tags             []Tags       `json:"tags" bson:"tags" validate:"nonnil,nonzero"`
	ListOfSteps      []string     `json:"listOfSteps" bson:"list_of_steps" validate:"nonnil,nonzero"`
	ListOfCategories []Category   `json:"listOfCategories" bson:"list_of_categories" validate:"nonnil,nonzero"`
	Description      string       `json:"description" bson:"description" validate:"nonnil,nonzero"`
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
func GetAllRecipes(client *mongo.Client, db string, page, limit string) ([]Recipe, int64, error) {
	recipes := []Recipe{}
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, err
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, 0, err
	}
	findOptions := options.Find()
	findOptions.SetLimit(int64(l))
	findOptions.SetSkip(int64((p - 1) * l))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionRecipes)
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(ctx) {
		var recipe Recipe
		if err = cursor.Decode(&recipe); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, 0, nil
}
func GetAllRecipesByTags(client *mongo.Client, db string, tags []int, page, limit string) ([]Recipe, int64, error) {
	fmt.Println("tags", tags)
	recipes := []Recipe{}
	recipe := Recipe{}
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, err
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionRecipes)
	findOptions := options.Find()
	findOptions.SetLimit(int64(l))
	findOptions.SetSkip(int64((p - 1) * l))
	cursor, err := collection.Find(ctx, bson.M{"tags": bson.M{"$in": tags}}, findOptions)
	for cursor.Next(ctx) {
		if err = cursor.Decode(&recipe); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}
	totalElements, err := collection.CountDocuments(ctx, bson.M{}, nil)
	if err != nil {
		return nil, 0, err
	}
	return recipes, totalElements, nil
}
