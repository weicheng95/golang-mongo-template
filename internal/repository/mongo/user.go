package mongo

import (
	"context"
	"time"

	"github.com/weicheng95/go-mongo-template/internal/constant/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(ctx context.Context, email string) (signedUser *model.User, err error)
	UpdateUserById(ctx context.Context, userId string, user model.User) error
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type User struct {
	collection *mongo.Collection
}

func UserInit(collection *mongo.Collection) *User {
	return &User{
		collection: collection,
	}
}

func (u *User) CreateUser(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().UTC().Unix()
	user.UpdatedAt = time.Now().UTC().Unix()
	user.UserId = user.ID.Hex()

	_, err := u.collection.InsertOne(ctx, *user)
	return err
}

func (u *User) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	//opt := options.FindOne().SetProjection(bson.M{"password": 0})
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	return user, err
}

func (u *User) UpdateUserById(ctx context.Context, userId string, user model.User) error {
	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	user.UpdatedAt = time.Now().UTC().Unix()
	us := &model.User{}
	return u.collection.FindOneAndUpdate(ctx, filter, user).Decode(us)
}

func (u *User) IsUserExist(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := u.collection.CountDocuments(ctx, filter)
	return count > 0, err
}
