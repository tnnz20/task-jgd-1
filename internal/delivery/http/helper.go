package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tnnz20/jgd-task-1/internal/model"
)

// GetIDFromPath extracts and parses an integer ID from the request path
func GetIDFromPath(r *http.Request, param string) (int, error) {
	idStr := r.PathValue(param)
	return strconv.Atoi(idStr)
}

// WriteJSON writes a JSON response to the client
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response to the client
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.WebResponse[any]{Errors: message})
}

// ReadJSON reads and decodes JSON from request body into the target
func ReadJSON(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}
