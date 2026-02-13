package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// PostsService holds dependencies for posts endpoints
type PostsService struct {
	DB            *sql.DB
	PostsStoreDir string // Path to .posts_store directory
}

// PostSubmitResponse is returned after a successful post submission
type PostSubmitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	PostID  int64  `json:"postId,omitempty"`
}

// ErrorResponse is returned when an error occurs
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Submit handles POST /api/posts/submit
// Accepts multipart form data with:
// - andrewId: string
// - name: string
// - caption: string (max 256 chars)
// - instagramUsername: string (Instagram handle without @)
// - file_0 through file_9: uploaded files (images/videos)
func (s *PostsService) Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form (max 100MB)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to parse form data: " + err.Error(),
		})
		return
	}

	// Extract form fields
	andrewId := strings.TrimSpace(r.FormValue("andrewId"))
	name := strings.TrimSpace(r.FormValue("name"))
	instagramUsername := strings.TrimSpace(r.FormValue("instagramUsername"))
	caption := strings.TrimSpace(r.FormValue("caption"))

	// Validation
	if andrewId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Andrew ID is required",
		})
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Name is required",
		})
		return
	}

	if instagramUsername == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Instagram username is required",
		})
		return
	}

	// Remove @ prefix if present
	if strings.HasPrefix(instagramUsername, "@") {
		instagramUsername = instagramUsername[1:]
	}

	if caption == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Caption is required",
		})
		return
	}

	if len(caption) > 256 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Caption exceeds 256 characters",
		})
		return
	}

	// Create unique directory for this submission
	// Format: .posts_store/[andrewid]/[timestamp]/
	timestamp := time.Now().Unix()
	postDir := filepath.Join(s.PostsStoreDir, andrewId, fmt.Sprintf("%d", timestamp))

	// Create directory structure
	err = os.MkdirAll(postDir, 0755)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to create post directory: " + err.Error(),
		})
		return
	}

	// Write caption to caption.txt
	captionPath := filepath.Join(postDir, "caption.txt")
	err = os.WriteFile(captionPath, []byte(caption), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to save Instagram username: " + err.Error(),
		})
		return
	}

	// Write Instagram username to instagram.txt
	instagramPath := filepath.Join(postDir, "instagram.txt")
	err = os.WriteFile(instagramPath, []byte(instagramUsername), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to save caption: " + err.Error(),
		})
		return
	}

	// Process uploaded files
	filesProcessed := 0
	for i := 0; i < 10; i++ {
		fileKey := fmt.Sprintf("file_%d", i)
		file, header, err := r.FormFile(fileKey)
		if err != nil {
			// No more files
			continue
		}
		defer file.Close()

		// Get file extension from original filename
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == "" {
			// Try to determine from content type
			contentType := header.Header.Get("Content-Type")
			switch {
			case strings.HasPrefix(contentType, "image/jpeg"):
				ext = ".jpg"
			case strings.HasPrefix(contentType, "image/png"):
				ext = ".png"
			case strings.HasPrefix(contentType, "image/gif"):
				ext = ".gif"
			case strings.HasPrefix(contentType, "image/webp"):
				ext = ".webp"
			case strings.HasPrefix(contentType, "video/mp4"):
				ext = ".mp4"
			case strings.HasPrefix(contentType, "video/quicktime"):
				ext = ".mov"
			case strings.HasPrefix(contentType, "video/webm"):
				ext = ".webm"
			default:
				ext = ".bin"
			}
		}

		// Validate file type
		validExtensions := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
			".mp4": true, ".mov": true, ".webm": true, ".avi": true,
		}
		if !validExtensions[ext] {
			// Skip invalid file types
			continue
		}

		// Save file with numeric name
		destPath := filepath.Join(postDir, fmt.Sprintf("%d%s", filesProcessed, ext))
		destFile, err := os.Create(destPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Success: false,
				Message: "Failed to save file: " + err.Error(),
			})
			return
		}

		_, err = io.Copy(destFile, file)
		destFile.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Success: false,
				Message: "Failed to write file: " + err.Error(),
			})
			return
		}

		filesProcessed++
	}

	if filesProcessed == 0 {
		// Clean up the directory since no files were uploaded
		os.RemoveAll(postDir)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "At least one image or video is required",
		})
		return
	}

	// Insert into database
	var postID int64
	err = s.DB.QueryRow(`
		INSERT INTO posts (andrewid, username, content, created_at, verified, approved)
		VALUES ($1, $2, $3, NOW(), false, false)
		RETURNING id
	`, andrewId, name, postDir).Scan(&postID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to save post to database: " + err.Error(),
		})
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PostSubmitResponse{
		Success: true,
		Message: fmt.Sprintf("Post submitted successfully with %d file(s)", filesProcessed),
		PostID:  postID,
	})
}

// List handles GET /api/posts/list
// Returns all posts for the authenticated user
func (s *PostsService) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	andrewId := r.URL.Query().Get("andrewId")
	if andrewId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Andrew ID is required",
		})
		return
	}

	rows, err := s.DB.Query(`
		SELECT id, andrewid, username, content, created_at, verified, approved
		FROM posts
		WHERE andrewid = $1
		ORDER BY created_at DESC
	`, andrewId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	type Post struct {
		ID        int64     `json:"id"`
		AndrewID  string    `json:"andrewId"`
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"createdAt"`
		Verified  bool      `json:"verified"`
		Approved  bool      `json:"approved"`
	}

	posts := []Post{}
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.AndrewID, &p.Username, &p.Content, &p.CreatedAt, &p.Verified, &p.Approved)
		if err != nil {
			continue
		}
		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"posts":   posts,
		"count":   len(posts),
	})
}

// GetPost handles GET /api/posts/{id}
// Returns a single post by ID
func (s *PostsService) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract post ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Post ID is required",
		})
		return
	}

	postIdStr := pathParts[len(pathParts)-1]
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	type PostDetail struct {
		ID        int64     `json:"id"`
		AndrewID  string    `json:"andrewId"`
		Username  string    `json:"username"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"createdAt"`
		Verified  bool      `json:"verified"`
		Approved  bool      `json:"approved"`
		Caption   string    `json:"caption"`
		Files     []string  `json:"files"`
	}

	var p PostDetail
	err = s.DB.QueryRow(`
		SELECT id, andrewid, username, content, created_at, verified, approved
		FROM posts
		WHERE id = $1
	`, postId).Scan(&p.ID, &p.AndrewID, &p.Username, &p.Content, &p.CreatedAt, &p.Verified, &p.Approved)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Post not found",
		})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Message: "Failed to fetch post: " + err.Error(),
		})
		return
	}

	// Read caption from file
	captionPath := filepath.Join(p.Content, "caption.txt")
	captionBytes, err := os.ReadFile(captionPath)
	if err == nil {
		p.Caption = string(captionBytes)
	}

	// List files in the post directory
	entries, err := os.ReadDir(p.Content)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != "caption.txt" {
				p.Files = append(p.Files, entry.Name())
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"post":    p,
	})
}
