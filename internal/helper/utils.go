package helper

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}

type apiFn func(http.ResponseWriter, *http.Request) error

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: data,
	})

}

func WriteError(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(APIError{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	})
}

func MakeHTTPHandlerFunc(fn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := fn(w, r); err != nil {
			if e, ok := err.(APIError); ok {
				slog.Error("API error:", "err", e, "status", e.Status)
				WriteJSON(w, e.Status, e)
			}
		}
	}
}

func PermissionError() APIError {
	return APIError{
		Status:  http.StatusForbidden,
		Message: "You do not have permission to access this resource",
	}
}
