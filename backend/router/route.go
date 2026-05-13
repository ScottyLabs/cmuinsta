// Package router handles the Gin routing task and registers endpoints to
// their handlers
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/scottylabs/cmuinsta/backend/handlers"
	"github.com/scottylabs/cmuinsta/backend/middleware"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// Authentication endpoint
	r.POST("/api/v1/auth/callback", handlers.AuthCallback)

	api := r.Group("/api/v1")
	api.Use(middleware.Authenticate())
	{
		api.POST("/me", handlers.CreateUser)
		api.GET("/me", handlers.GetUser)
	}

	return r
}
