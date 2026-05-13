// Package middleware includes the middleware for the Gin REST API
//
// It takes the JWT required for all endpoints, verifies it against the
// ScottyLabs Keycloak OIDC provider, and injects the corresponding
// Andrew ID into the request context for handlers to use.
package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

var verifier *oidc.IDTokenVerifier

// Init fetches Keycloak's public keys and initializes the token verifier.
// Must be called once at startup in main.go before the server starts.
func Init(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, os.Getenv("OIDC_ISSUER_URL"))
	if err != nil {
		return err
	}
	verifier = provider.Verifier(&oidc.Config{
		ClientID: os.Getenv("OIDC_CLIENT_ID"),
	})
	return nil
}

// Authenticate injects the Andrew ID into a gin.Context object by verifying
// the provided JWT against the Keycloak OIDC provider. It aborts with an
// error when the JWT is missing, malformed, expired, or invalid.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the auth token
		authHeader := c.GetHeader("Authorization")
		// If empty or not Bearer auth then we reject
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or malformed Authorization header",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify signature, expiry, issuer, and audience against Keycloak's
		// public keys — fetched automatically via OIDC discovery on startup.
		idToken, err := verifier.Verify(c.Request.Context(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired JWT",
			})
			return
		}

		// Extract claims from the verified token
		var claims struct {
			PreferredUsername string `json:"preferred_username"`
		}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid JWT claims",
			})
			return
		}

		// Extract preferred_username — this IS the andrewid.
		// Never trust the request body for this.
		if claims.PreferredUsername == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "JWT missing preferred_username claim",
			})
			return
		}

		// Inject into context. Handlers retrieve this with c.GetString("andrewid").
		c.Set("andrewid", claims.PreferredUsername)
		// Pass to the next handler.
		c.Next()
	}
}
