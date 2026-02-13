package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// PrefroshService holds dependencies for prefrosh endpoints
type PrefroshService struct {
	DB *sql.DB
}

// List returns a list of data (formerly /api/data)
// GET /api/prefrosh/list
func (s *PrefroshService) List(w http.ResponseWriter, r *http.Request) {
	// Example: Query database here in the future
	response := map[string]string{
		"message": "Hello from Go! (Prefrosh Service)",
		"time":    "Now",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
