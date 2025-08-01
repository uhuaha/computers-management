package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bdlm/log"
)

// handleError writes a JSON-formatted error response with the given message and HTTP status code.
// It sets the "Content-Type" header to "application/json" and encodes the error as:
// {"error": "<errMsg>"}.
func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": errMsg,
	}); err != nil {
		log.Error("failed to encode error message: " + err.Error())
	}
}
