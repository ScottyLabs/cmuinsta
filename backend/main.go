package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
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
	// |------------|------------|------------|------------|-----------|-------------|
	// | id         | andrewid   | content    | created_at | approved  | approved_at |
	// |------------|------------|------------|------------|-----------|-------------|
	//
	// `content` contains a path to a folder, this folder will look something like this:
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
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			approved BOOLEAN DEFAULT FALSE,
			approved_at TIMESTAMP DEFAULT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("✅ Database schema initialized!")

	// 3. Setup Router
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Go backend is running"})
	})

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Example: Query database here in the future
		response := map[string]string{
			"message": "Hello from Go!",
			"time":    "Now",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// 4. CORS Middleware (for development)
	handler := corsMiddleware(mux)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
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
