// Package router handles the Gin routing task and registers endpoints to
// their handlers
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/scottylabs/cmuinsta/backend/handlers"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// Authentication endpoint
	r.POST("/api/v1/auth/callback", handlers.AuthCallback)
	r.POST("/api/v1/hook", handlers.DMHook)

	return r
}
