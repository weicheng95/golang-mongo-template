package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/weicheng95/go-mongo-template/internal/constant/model"
	"github.com/weicheng95/go-mongo-template/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	log = logger.NewLogger("go-mango-server")
)


func (sv *service) CreateUser(ctx context.Context, user *model.User) error {
	if err := sv.userRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (sv *service) Login(ctx context.Context, user *model.User) error {
	return nil
}

func (sv *service) UpdateUserById(ctx context.Context, userId string, user model.User) error {
	return nil
}

func (sv *service) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := sv.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("unable to find user")
		}
		return nil, err
	}
	return user, nil
}

func (sv *service) ValidateUserEmail(ctx context.Context, email string) ( err error) {
	valid, err := sv.userRepo.IsUserExist(ctx, email)

	if err != nil {
		log.Error(err)
	}

	if valid {
		err = fmt.Errorf("user existed")
	}
	return
}

//HashPassword is used to encrypt the password before it is stored in the DB
func (sv *service) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func (sv *service) VerifyPassword(providedPassword string, userPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return errors.New("login or password is incorrect")
	}
	return nil
}
