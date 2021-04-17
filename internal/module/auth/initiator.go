package auth

import (
	"context"
	"github.com/weicheng95/go-mongo-template/internal/constant/model"
	"github.com/weicheng95/go-mongo-template/internal/repository/mongo"
)

// contains the function of business logic of domain user
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	ValidateUserEmail(ctx context.Context, email string) (err error)
	HashPassword(password string) string
	VerifyPassword(providedPassword string, userPassword string) error
	FindUserByEmail(ctx context.Context, email string) (user *model.User, err error)
	UpdateUserById(ctx context.Context, userId string, user model.User) error
}

// port to connect service & repo / persistent / external
type service struct {
	userRepo       mongo.UserRepository
}

// Initialize takes all necessary service for domain user to run the business logic of domain user
func Initialize(
	userRepo mongo.UserRepository,
) UserService {
	return &service{
		userRepo,
	}
}