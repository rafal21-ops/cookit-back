package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Daniorocket/cookit-back/lib"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CollectionCategory = "categories"
var ErrCreateCategory = errors.New("Failed to create category record")
var ErrFindCategory = errors.New("Failed to find category record")

type Category struct {
	ID      string   `json:"id" bson:"id"`
	LabelPL string   `json:"labelPL" bson:"label_pl"`
	LabelEN string   `json:"labelEN" bson:"label_en"`
	File    lib.File `json:"file" bson:"file"`
}

func GetAllCategories(client *mongo.Client, db string) ([]Category, error) {
	categories := []Category{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionCategory)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var category Category
		if err = cursor.Decode(&category); err != nil {
			return nil, err
		}
		fmt.Println("Category:", category)
		categories = append(categories, category)
	}
	fmt.Println("Categories:", categories)
	return categories, nil
}
func CreateCategory(client *mongo.Client, db string, category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionCategory)
	if _, err := collection.InsertOne(ctx, &category); err != nil {
		return ErrCreateCategory
	}
	return nil
}
