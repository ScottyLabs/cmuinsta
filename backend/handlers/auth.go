package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthCallback handles the exchange of a code from the frontend OIDC flow to
// a JWT from Keycloak handled privately in the backend. This is in accordance
// with RFC 4 for authentication and RFC 2 POST `/auth/callback`
func AuthCallback(c *gin.Context) {
	var body struct {
		Code     string `json:"code" binding:"required"`
		Verifier string `json:"code_verifier" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or verifier"})
		return
	}

	resp, err := http.PostForm(
		os.Getenv("OIDC_ISSUER_URL")+"/protocol/openid-connect/token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {os.Getenv("OIDC_CLIENT_ID")},
			"redirect_uri":  {os.Getenv("OIDC_REDIRECT_URI")},
			"code":          {body.Code},
			"code_verifier": {body.Verifier},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact Keycloak"})
		return
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code exchange failed"})
		return
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	// Return the JWT to the frontend
	c.JSON(http.StatusOK, gin.H{
		"token":      tokenResponse.IDToken,
		"expires_in": tokenResponse.ExpiresIn,
	})
}
