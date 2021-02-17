package routing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Daniorocket/cookit-back/handlers"
	"github.com/Daniorocket/cookit-back/logger"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRouter() (*mux.Router, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		os.Getenv("MONGODB_URI"),
	))
	if err != nil {
		log.Println("Err connect mongo:", err)
		return nil, err
	}
	handler := handlers.Handler{
		Client:       client,
		DatabaseName: "CookIt",
	}
	collection := client.Database("CookIt").Collection("users")
	//Index for users
	keysUser := []string{"id", "email", "username"}
	for i := range keysUser {
		indexName, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.M{
					keysUser[i]: 1,
				},
				Options: options.Index().SetUnique(true).SetBackground(true).SetSparse(true),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(indexName)
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range initRoutes(handler) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router, nil
}
