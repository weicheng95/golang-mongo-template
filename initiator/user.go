package initiator

import (
	"github.com/weicheng95/go-mongo-template/internal/handler/rest"
	"github.com/weicheng95/go-mongo-template/internal/module/auth"
	mongodb "github.com/weicheng95/go-mongo-template/internal/repository/mongo"
	"github.com/weicheng95/go-mongo-template/pkg/helper"
	"go.mongodb.org/mongo-driver/mongo"
)

// User initializes the domain user
func UserRestInit(client *mongo.Client) *rest.UserHandler {
	userCollection := helper.GetMongoDBCollection(client, "sampledb", "user")

	userRepo := mongodb.UserInit(userCollection)

	userService := auth.Initialize(userRepo)

	handler := rest.UserHandlerInit(userService)

	return handler
}
