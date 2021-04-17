package rest

import (
	"context"
	"encoding/json"
	"net/http"
)

type Error struct {
	Msg  string `json:"msg"`
}

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error Error `json:"error"`
}

func RenderErrorResponse(ctx context.Context, w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	resp := ErrorResponse{Error: Error{Msg: err.Error()}}
	RenderResponse(w, resp, status)
}

func RenderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		// XXX Do something with the error ;)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// XXX Do something with the error ;)
	}
}
