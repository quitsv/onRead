package Controllers

import (
	"encoding/json"
	"net/http"
)

func PrintSuccess(status int, message string, w http.ResponseWriter) {
	response := SuccessResponse{Status: status, Message: message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PrintError(status int, message string, w http.ResponseWriter) {
	response := ErrorResponse{Status: status, Message: message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
