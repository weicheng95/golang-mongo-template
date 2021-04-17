package model

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

//User is the port between module and repo
type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	FirstName *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName  *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password  *string            `json:"password,omitempty" validate:"required,min=6"`
	Email     *string            `json:"email" validate:"email,required"`
	Phone     *string            `json:"phone" validate:"required"`
	CreatedAt int64              `json:"created_at, omitempty"`
	UpdatedAt int64              `json:"updated_at, omitempty"`
	UserId    string             `json:"user_id"`
	Role      string             `json:"role"`
}

type SignupUserRequest struct {
	*User
	Role string `json:"role"`
}

func (a *SignupUserRequest) Bind(r *http.Request) error {
	if a.User == nil || a.Email == nil || a.FirstName == nil || a.LastName == nil || a.Password == nil {
		return errors.New("missing required fields")
	}
	return nil
}

type LoginUserRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (a *LoginUserRequest) Bind(r *http.Request) error {
	if a.Email == nil || a.Password == nil {
		return errors.New("missing required fields")
	}
	return nil
}

type LoginUserResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredIn    int64  `json:"expired_in"`
}

type UserClaims struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	jwt.StandardClaims
}

//func (u *User) validate() url.Values {
//	errs := url.Values{}
//
//	// check if the title empty
//	if a.Title == "" {
//		errs.Add("title", "The title field is required!")
//	}
//
//	// check the title field is between 3 to 120 chars
//	if len(a.Title) < 3 || len(a.Title) > 120 {
//		errs.Add("title", "The title field must be between 3-120 chars!")
//	}
//
//	if a.Body == "" {
//		errs.Add("body", "The body field is required!")
//	}
//
//	if len(a.Body) < 50 || len(a.Body) > 500 {
//		errs.Add("body", "The title field must be between 50-500 chars!")
//	}
//
//	// Note: for checking the complex or regex related validation I use govalidator here
//
//	return errs
//}
