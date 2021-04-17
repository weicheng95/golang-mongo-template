package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient() (*mongo.Client) {
	get := func(v string) string {
		res, err := Get(v)
		if err != nil {
			log.Fatalf("%s", err)
		}

		return res
	}

	mongodbHOST := get("database.mongodb_host")
	// mongodbDATABASE := get("database.mongodb_database")
	mongodbUsername := get("database.mongodb_username")
	mongodbPassword := get("database.mongodb_password")

	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", mongodbUsername, mongodbPassword, mongodbHOST)
	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Exit error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Exit error: %v", err)
	}

	return client
}

func GetMongoDBCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}
