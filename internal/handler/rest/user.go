package rest

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/weicheng95/go-mongo-template/internal/common/server"
	"github.com/weicheng95/go-mongo-template/internal/common/server/httperr"
	"github.com/weicheng95/go-mongo-template/internal/constant/model"
	"github.com/weicheng95/go-mongo-template/internal/module/auth"
	"net/http"

	"github.com/weicheng95/go-mongo-template/pkg/rest"
)

type UserHandler struct {
	userService auth.UserService
}

type RegisterResponse struct {
	Email string `json:"email"`
}

func UserHandlerInit(userService auth.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register connects the handlers to the router.
func (a *UserHandler) AuthRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/signup", a.signup)
	r.Post("/login", a.login)
	return r
}

func (a *UserHandler) UserRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/info", a.userInfo)
	return r
}


func (a *UserHandler) signup(w http.ResponseWriter, r *http.Request) {
	// deserialize the input, validate payload
	signupReq := &model.SignupUserRequest{}
	if err := render.Bind(r, signupReq); err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// validate user email
	if err := a.userService.ValidateUserEmail(r.Context(), *signupReq.Email); err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// hash password
	hashPassword := a.userService.HashPassword(*signupReq.Password)
	if hashPassword == "" {
		rest.RenderErrorResponse(r.Context(), w, http.StatusInternalServerError, fmt.Errorf("something wrong"))
		return
	}
	signupReq.Password = &hashPassword

	// create user
	if err := a.userService.CreateUser(r.Context(), signupReq.User); err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// generate auth token
	token, refreshToken, expiredIn, err := auth.GenerateToken(model.UserClaims{
		FirstName: *signupReq.FirstName,
		LastName: *signupReq.LastName,
		Email: *signupReq.Email,
		UserId: signupReq.UserId,
	})

	if err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	signedUser := model.LoginUserResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiredIn:    expiredIn,
	}

	fmt.Println("hello")

	render.Respond(w, r, signedUser)
	render.Status(r, http.StatusOK)
}

func (a * UserHandler) login(w http.ResponseWriter, r *http.Request) {
	// validate payload, headers (auth token)

	// deserialize the input
	loginReq := &model.LoginUserRequest{}
	if err := render.Bind(r, loginReq); err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// find user
	userInfo, err := a.userService.FindUserByEmail(r.Context(), *loginReq.Email)
	if err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// compare password
	if err := a.userService.VerifyPassword(*userInfo.Password, *loginReq.Password); err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// generate auth token
	token, refreshToken, expiredIn, err := auth.GenerateToken(model.UserClaims{
		FirstName: *userInfo.FirstName,
		LastName: *userInfo.LastName,
		Email: *userInfo.Email,
		UserId: userInfo.UserId,
	})

	if err != nil {
		rest.RenderErrorResponse(r.Context(), w, http.StatusInternalServerError, errors.New("something wrong"))
	}

	signedUser := model.LoginUserResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiredIn:    expiredIn,
	}
	render.Respond(w, r, signedUser)
	render.Status(r, http.StatusOK)
}

func (a *UserHandler) userInfo(w http.ResponseWriter, r *http.Request) {
	// get user from ctx
	u, err := server.UserFromCtx(r.Context())

	if err != nil {
		httperr.Unauthorised(err, w, r)
		return
	}

	userInfo, err := a.userService.FindUserByEmail(r.Context(), u.Email)
	render.Respond(w, r, userInfo)
	render.Status(r, http.StatusOK)
}