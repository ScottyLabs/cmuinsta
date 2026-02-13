package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// VerificationService holds dependencies for verification endpoints
type VerificationService struct {
	DB *sql.DB
}

// List returns a list of verification data
// GET /api/verify/list
func (s *VerificationService) List(w http.ResponseWriter, r *http.Request) {
	// Example: Query database here in the future
	response := map[string]string{
		"message": "Hello from Go! (Verification Service)",
		"time":    "Now",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
