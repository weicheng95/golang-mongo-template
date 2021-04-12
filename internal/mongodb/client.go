package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/weicheng95/go-mongo-template/internal/envvar"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB() (*mongo.Client, error) {
	conf := &envvar.Configuration{}
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			log.Fatalf("Couldn't get configuration value for %s: %s", v, err)
		}

		return res
	}

	mongodbHOST := get("MONGODB_HOST")
	mongodbDATABASE := get("MONGODB_HOST")
	mongodbUsername := get("MONGODB_USERNAME")
	mongodbPassword := get("MONGODB_PASSWORD")

	var ctx = context.Background()
	// "mongodb+srv://admin:admin@cluster0.h6dcj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", mongodbUsername, mongodbPassword, mongodbHOST, mongodbDATABASE)
	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
