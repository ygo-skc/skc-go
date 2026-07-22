package model

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// =======================
// Error Response
// =======================
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string { return e.Message }

func (e *APIError) HandleServerResponse(res http.ResponseWriter) {
	if e.StatusCode == 0 {
		e.StatusCode = 500 // default error code
	}
	res.WriteHeader(e.StatusCode)
	if err := json.NewEncoder(res).Encode(e); err != nil {
		slog.Error("Failed to encode API error response", slog.Any("err", err))
	}
}

func HandleServerResponse(apiErr APIError, res http.ResponseWriter) {
	apiErr.HandleServerResponse(res)
}

// =======================
// Success Response
// =======================
type Success struct {
	Message string `json:"message"`
}
