package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJson   = "application/json"
)

type APIError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, message string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func JsonResponse(body interface{}, w http.ResponseWriter) {
	w.Header().Set(contentTypeHeader, contentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func InvalidJsonResponse(w http.ResponseWriter) {
	apiError := NewAPIError(http.StatusBadRequest, "Invalid json body")

	w.Header().Set(contentTypeHeader, contentTypeJson)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(apiError)
}

func ErrorJsonResponse(w http.ResponseWriter, message string) {
	apiError := NewAPIError(http.StatusBadRequest, message)

	w.Header().Set(contentTypeHeader, contentTypeJson)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(apiError)
}
