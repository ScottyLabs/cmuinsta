package handlers

import (
	"database/sql"
	"net/http"
)

// AdminService holds dependencies for admin endpoints
type AdminService struct {
	DB *sql.DB
}

// Dashboard is a sample admin endpoint
// GET /api/admin/dashboard
func (s *AdminService) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Admin Dashboard"))
}
