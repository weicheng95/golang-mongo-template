package mongo_test

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mongodb "github.com/weicheng95/go-mongo-template/internal/repository/mongo"
	"github.com/weicheng95/go-mongo-template/pkg/helper"
	"testing"
)

func TestRepository(t *testing.T) {
	var env, configFile string
	flag.StringVar(&env, "env", "development", "Environment Variables filename")
	flag.StringVar(&configFile, "configFile", "../../../config.yaml", "Environment Variables filename")
	flag.Parse()

	err := helper.Load(env, configFile);
	require.NoError(t, err)
	dbClient := helper.NewMongoDBClient()
	require.NoError(t, err)
	userCollection := helper.GetMongoDBCollection(dbClient, "sampledb", "user")
	userRepo := mongodb.UserInit(userCollection)
	testFindUserByEmail(t, userRepo)
	testUserExist(t, userRepo)
}
func testFindUserByEmail(t *testing.T, userRepo *mongodb.User) {
	user, err := userRepo.FindUserByEmail(context.Background(), "day15user1@gmail.com")
	require.NoError(t, err)
	require.NotEmpty(t, user)
	assert.Equal(t, *user.Email, "day15user1@gmail.com")
	assert.Nil(t, user.Password)
}

func testUserExist(t *testing.T, userRepo *mongodb.User) {
	valid, err := userRepo.IsUserExist(context.Background(), "day15user1@gmail.com")
	require.NoError(t, err)
	assert.Equal(t, valid, true)
}
