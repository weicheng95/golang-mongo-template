package main

//import (
//	"fmt"
//	"github.com/weicheng95/go-mongo-template/internal/auth"
//	"github.com/weicheng95/go-mongo-template/internal/auth/models"
//)
//
//func main() {
//	claims := models.SignedDetails{
//		Email: "jerry@gmail.com",
//		FirstName: "Jerry",
//		LastName: "chew",
//		UserId: "1231241161",
//	}
//	signedToken, signedRefreshToken, err := auth.GenerateTokens(claims)
//
//	if err != nil {
//		fmt.Errorf("%w", err)
//	}
//
//	fmt.Println(signedToken)
//	fmt.Println(signedRefreshToken)
//
//
//	originalClaims, err := auth.ValidateToken(signedToken)
//	if err != nil {
//		fmt.Errorf("%w", err)
//	}
//
//	fmt.Printf("%+v", originalClaims)
//}


