package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseWithJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func ResponseWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Status:  http.StatusText(statusCode),
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func ResponseWithSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Status:  http.StatusText(statusCode),
		Data: data,
	}

	json.NewEncoder(w).Encode(response)
}
