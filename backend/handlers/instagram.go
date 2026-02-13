package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// InstagramService holds dependencies for Instagram-related endpoints
type InstagramService struct{}

// ValidateResponse is returned when validating an Instagram username
type ValidateResponse struct {
	Success  bool   `json:"success"`
	Username string `json:"username"`
	Exists   bool   `json:"exists"`
	Message  string `json:"message,omitempty"`
}

// Validate checks if an Instagram username exists
// GET /api/instagram/validate?username=xxx
func (s *InstagramService) Validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidateResponse{
			Success: false,
			Message: "Username is required",
		})
		return
	}

	// Validate username format (alphanumeric, underscores, periods, 1-30 chars)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9._]{1,30}$`)
	if !usernameRegex.MatchString(username) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidateResponse{
			Success:  false,
			Username: username,
			Exists:   false,
			Message:  "Invalid username format",
		})
		return
	}

	// Check if the Instagram profile exists by making a request to Instagram
	exists, err := checkInstagramProfileExists(username)
	if err != nil {
		// If we can't verify, we'll return success but with exists=true to not block the user
		// This is a graceful degradation - we don't want to block users if Instagram is down
		json.NewEncoder(w).Encode(ValidateResponse{
			Success:  true,
			Username: username,
			Exists:   true, // Assume exists if we can't verify
			Message:  "Could not verify, assuming valid",
		})
		return
	}

	json.NewEncoder(w).Encode(ValidateResponse{
		Success:  true,
		Username: username,
		Exists:   exists,
	})
}

// checkInstagramProfileExists makes a HEAD request to Instagram to check if a profile exists
func checkInstagramProfileExists(username string) (bool, error) {
	// Instagram profile URL
	profileURL := fmt.Sprintf("https://www.instagram.com/%s/", username)

	// Create a client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
		// Don't follow redirects - we want to check the initial response
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Make a HEAD request to check if the profile exists
	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return false, err
	}

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Instagram returns 200 for existing profiles, 404 for non-existing ones
	// However, Instagram may also return 200 with a "Page Not Found" content for some cases
	// A 200 status code generally means the profile exists
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	// 404 means the profile doesn't exist
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	// For other status codes (like 429 rate limit), we can't be sure
	// Return an error to trigger the graceful degradation
	if resp.StatusCode == http.StatusTooManyRequests {
		return false, fmt.Errorf("rate limited by Instagram")
	}

	// For redirects to login page, the profile might be private but exists
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
