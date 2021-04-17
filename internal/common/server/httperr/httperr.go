package httperr

import (
	"github.com/go-chi/render"
	"github.com/weicheng95/go-mongo-template/internal/common/server/logging"
	"net/http"
)

func InternalError(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(w, r, err, "Internal server error", http.StatusInternalServerError)
}

func Unauthorised(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(w, r, err, "Unauthorised", http.StatusUnauthorized)
}

func BadRequest(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(w, r, err, "Bad request", http.StatusBadRequest)
}

func httpRespondWithError(w http.ResponseWriter, r *http.Request, err error, msg string, status int) {
	logger := logging.GetLogEntry(r)
	logger.WithError(err).Warn(msg)
	resp := ErrorResponse{msg, status}
	w.Header().Set("Content-Type", "application/json")
	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Message    string `json:"message"`
	HttpStatus int `json:"status"`
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.HttpStatus)
	return nil
}
