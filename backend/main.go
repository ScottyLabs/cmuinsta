package main

import (
	"cmuinsta/backend/handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	// 0. Debug OIDC Configuration
	fmt.Println("🔐 OIDC Configuration:")
	fmt.Printf("   OIDC_ISSUER_URL: %s\n", os.Getenv("OIDC_ISSUER_URL"))
	fmt.Printf("   OIDC_CLIENT_ID: %s\n", os.Getenv("OIDC_CLIENT_ID"))
	fmt.Printf("   OIDC_CLIENT_SECRET: %s\n", maskSecret(os.Getenv("OIDC_CLIENT_SECRET")))
	fmt.Printf("   OIDC_REDIRECT_URI: %s\n", os.Getenv("OIDC_REDIRECT_URI"))
	fmt.Printf("   ADMIN_IDS: %s\n", os.Getenv("ADMIN_IDS"))

	// 1. Database Connection
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("DEBUG: The URL Go is using is: %s\n", os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("✅ Connected to Postgres database!")
	var dbName string
	err = db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err == nil {
		fmt.Printf("📊 Confirmed: Connected to database [%s]\n", dbName)
	}

	// 2. Create Table
	// The table format is like this
	// |------------|------------|------------|------------|--------------|-------------|-----------|
	// | id         | andrewid   | username   | created_at | submitted_at | approved_at | posted_at |
	// |------------|------------|------------|------------|--------------|-------------|-----------|
	//
	// content folder will look something like this:
	//			 andrewid
	//			  ├── caption.txt    (required)
	//			  ├── 0.jpg          (required)
	//			  ├── 1.mov          (optional)
	//			  │   ...            (optional)
	//			  └── 9.mp4          (optional)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			andrewid VARCHAR(8) NOT NULL,
			username VARCHAR(30) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			submitted_at TIMESTAMP DEFAULT NULL,
			approved_at TIMESTAMP DEFAULT NULL,
			posted_at TIMESTAMP DEFAULT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("✅ Database schema initialized!")

	// 3. Setup posts store directory
	postsStoreDir := os.Getenv("POSTS_STORE_DIR")
	if postsStoreDir == "" {
		// Default to .posts_store in the project root (one level up from backend/)
		postsStoreDir = "../.posts_store"
	}
	// Ensure the directory exists
	if err := os.MkdirAll(postsStoreDir, 0755); err != nil {
		log.Fatalf("Error creating posts store directory: %v", err)
	}
	absPostsStoreDir, _ := filepath.Abs(postsStoreDir)
	fmt.Printf("📁 Posts store directory: %s\n", absPostsStoreDir)

	// 4. Initialize Handlers
	admin := &handlers.AdminService{DB: db}
	prefrosh := &handlers.PrefroshService{DB: db}
	auth := &handlers.AuthService{DB: db}
	posts := &handlers.PostsService{DB: db, PostsStoreDir: absPostsStoreDir}
	instagram := &handlers.InstagramService{}

	// 5. Setup Router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Go backend is running"})
	})

	// Auth routes
	mux.HandleFunc("POST /api/auth/callback", auth.Callback)
	mux.HandleFunc("GET /api/auth/login-url", auth.GetLoginURL)
	mux.HandleFunc("GET /api/auth/logout-url", auth.GetLogoutURL)
	mux.HandleFunc("POST /api/auth/check-admin", auth.CheckAdmin)
	mux.HandleFunc("GET /api/auth/me", auth.UserInfo)

	// Posts routes
	mux.HandleFunc("POST /api/posts/submit", posts.Submit)
	mux.HandleFunc("GET /api/posts/list", posts.List)
	mux.HandleFunc("GET /api/posts/{id}", posts.GetPost)

	// Admin routes
	mux.HandleFunc("GET /api/admin/dashboard", admin.Dashboard)

	// Instagram routes
	mux.HandleFunc("GET /api/instagram/validate", instagram.Validate)

	// Prefrosh routes (legacy)
	mux.HandleFunc("GET /api/prefrosh/list", prefrosh.List)

	// 6. CORS Middleware (for development)
	handler := corsMiddleware(mux)

	// 7. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// maskSecret hides most of a secret string for safe logging
func maskSecret(s string) string {
	if len(s) == 0 {
		return "(not set)"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from the frontend (adjust origin as needed for prod)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
