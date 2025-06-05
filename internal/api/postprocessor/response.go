package postprocessor

import (
	"encoding/json"
	"net/http"
)

// FormatResponse formats the response data before sending it to the client.
func FormatResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
