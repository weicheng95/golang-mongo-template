package server

import (
	"context"
	"fmt"
	"github.com/weicheng95/go-mongo-template/internal/common/server/httperr"
	"github.com/weicheng95/go-mongo-template/internal/constant/model"
	"github.com/weicheng95/go-mongo-template/internal/module/auth"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

func (a AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken := a.tokenFromHeader(r)
		if bearerToken == "" {
			httperr.Unauthorised(nil, w, r)
			return
		}

		token, err := auth.ValidateToken(bearerToken)
		if err != nil {
			httperr.Unauthorised(err, w, r)
			return
		}

		// it's always a good idea to use custom type as context value (in this case ctxKey)
		// because nobody from the outside of the package will be able to override/read this value
		ctx = context.WithValue(ctx, userContextKey, &model.UserClaims{
			UserId: token.UserId,
			Email:  token.Email,
			LastName: token.LastName,
			FirstName: token.FirstName,
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")
	// split token out from bearer (bearer ...)
	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}
	return ""
}

func UserFromCtx(ctx context.Context) (*model.UserClaims, error) {
	u, ok := ctx.Value(userContextKey).(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return u, nil
}
