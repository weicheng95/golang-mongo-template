package auth

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/weicheng95/go-mongo-template/internal/constant/model"
)


var (
	secretKey string = os.Getenv("auth.secret_key")
)

func GenerateToken(claims model.UserClaims) (token string, refreshToken string, expiredIn int64, err error) {

	expiredIn = time.Now().Local().AddDate(1, 0, 0).Unix()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expiredIn, // add 1 year
	}

	refreshClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().AddDate(1, 6, 0).Unix(), // add 1.6 years
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", 0, err
	}
	return
}

//ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *model.UserClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("the token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

//UpdateAllTokens renews the user tokens when they login
// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var updateObj primitive.D

// 	updateObj = append(updateObj, bson.E{"token", signedToken})
// 	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

// 	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

// 	upsert := true
// 	filter := bson.M{"user_id": userId}
// 	opt := options.UpdateOptions{
// 		Upsert: &upsert,
// 	}

// 	_, err := userCollection.UpdateOne(
// 		ctx,
// 		filter,
// 		bson.D{
// 			{"$set", updateObj},
// 		},
// 		&opt,
// 	)
// 	defer cancel()

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}

// 	return
// }
