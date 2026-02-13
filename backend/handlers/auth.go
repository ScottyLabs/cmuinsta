package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// AuthService holds dependencies for auth endpoints
type AuthService struct {
	DB *sql.DB
}

// OIDCConfig holds OIDC configuration
type OIDCConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// TokenResponse from Keycloak
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// UserInfo from Keycloak userinfo endpoint
type UserInfoResponse struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"` // Andrew ID
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Name              string `json:"name"`
}

// CallbackRequest from frontend
type CallbackRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirectUri,omitempty"`
}

// AuthResponse returned to frontend after successful auth
type AuthResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	User        struct {
		AndrewID  string `json:"andrewId"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		GivenName string `json:"givenName"`
	} `json:"user"`
	IsAdmin bool `json:"isAdmin"`
}

// getOIDCConfig returns OIDC configuration from environment
func getOIDCConfig() OIDCConfig {
	return OIDCConfig{
		IssuerURL:    os.Getenv("OIDC_ISSUER_URL"),
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("OIDC_REDIRECT_URI"),
	}
}

// Callback handles the OAuth callback - exchanges code for tokens
// POST /api/auth/callback
func (s *AuthService) Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Authorization code is required",
		})
		return
	}

	config := getOIDCConfig()

	// Use provided redirect URI or fall back to config
	redirectURI := req.RedirectURI
	if redirectURI == "" {
		redirectURI = config.RedirectURI
	}

	// Exchange code for tokens
	tokenURL := fmt.Sprintf("%s/protocol/openid-connect/token", config.IssuerURL)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", req.Code)
	data.Set("redirect_uri", redirectURI)

	client := &http.Client{Timeout: 10 * time.Second}
	tokenResp, err := client.PostForm(tokenURL, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to exchange code for tokens",
		})
		return
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(tokenResp.Body)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Token exchange failed",
			"details": string(body),
		})
		return
	}

	var tokens TokenResponse
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokens); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse token response",
		})
		return
	}

	// Fetch user info using the access token
	userInfoURL := fmt.Sprintf("%s/protocol/openid-connect/userinfo", config.IssuerURL)
	userInfoReq, _ := http.NewRequest("GET", userInfoURL, nil)
	userInfoReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch user info",
		})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo UserInfoResponse
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse user info",
		})
		return
	}

	// Check if user is admin
	andrewID := strings.ToLower(strings.TrimSpace(userInfo.PreferredUsername))
	isAdmin := checkIsAdmin(andrewID)

	// Build response
	response := AuthResponse{
		Success:     true,
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
		IsAdmin:     isAdmin,
	}
	response.User.AndrewID = andrewID
	response.User.Email = userInfo.Email
	response.User.Name = userInfo.Name
	response.User.GivenName = userInfo.GivenName

	json.NewEncoder(w).Encode(response)
}

// CheckAdmin checks if the given Andrew ID is in the ADMIN_IDS environment variable
// POST /api/auth/check-admin
func (s *AuthService) CheckAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		AndrewID string `json:"andrewId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.AndrewID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "andrewId is required",
		})
		return
	}

	andrewID := strings.ToLower(strings.TrimSpace(req.AndrewID))
	isAdmin := checkIsAdmin(andrewID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"isAdmin":  isAdmin,
		"andrewId": andrewID,
	})
}

// GetLoginURL returns the OIDC authorization URL
// GET /api/auth/login-url
func (s *AuthService) GetLoginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	config := getOIDCConfig()

	// Build authorization URL
	authURL := fmt.Sprintf("%s/protocol/openid-connect/auth", config.IssuerURL)

	params := url.Values{}
	params.Set("client_id", config.ClientID)
	params.Set("response_type", "code")
	params.Set("scope", "openid profile email")
	params.Set("redirect_uri", config.RedirectURI)

	loginURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"loginUrl": loginURL,
	})
}

// Logout handles logout by redirecting to Keycloak logout
// GET /api/auth/logout-url
func (s *AuthService) GetLogoutURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	config := getOIDCConfig()

	// Build logout URL
	logoutURL := fmt.Sprintf("%s/protocol/openid-connect/logout", config.IssuerURL)

	params := url.Values{}
	params.Set("client_id", config.ClientID)
	// Redirect to home after logout
	params.Set("post_logout_redirect_uri", strings.TrimSuffix(config.RedirectURI, "/callback"))

	fullLogoutURL := fmt.Sprintf("%s?%s", logoutURL, params.Encode())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"logoutUrl": fullLogoutURL,
	})
}

// UserInfo returns user information from a valid access token
// GET /api/auth/me
func (s *AuthService) UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	config := getOIDCConfig()

	// Validate token by calling userinfo endpoint
	userInfoURL := fmt.Sprintf("%s/protocol/openid-connect/userinfo", config.IssuerURL)
	client := &http.Client{Timeout: 10 * time.Second}

	userInfoReq, _ := http.NewRequest("GET", userInfoURL, nil)
	userInfoReq.Header.Set("Authorization", "Bearer "+accessToken)

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to validate token",
		})
		return
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid or expired token",
		})
		return
	}

	var userInfo UserInfoResponse
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse user info",
		})
		return
	}

	andrewID := strings.ToLower(strings.TrimSpace(userInfo.PreferredUsername))
	isAdmin := checkIsAdmin(andrewID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"andrewId":  andrewID,
			"email":     userInfo.Email,
			"name":      userInfo.Name,
			"givenName": userInfo.GivenName,
		},
		"isAdmin": isAdmin,
	})
}

// checkIsAdmin checks if an Andrew ID is in the admin list
func checkIsAdmin(andrewID string) bool {
	adminIDsEnv := os.Getenv("ADMIN_IDS")
	if adminIDsEnv == "" {
		return false
	}

	adminIDs := strings.Split(adminIDsEnv, ",")
	for _, id := range adminIDs {
		if strings.ToLower(strings.TrimSpace(id)) == andrewID {
			return true
		}
	}
	return false
}
