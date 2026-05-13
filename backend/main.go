package main

import (
	"context"
	"log"
	"os"

	"github.com/scottylabs/cmuinsta/backend/db"
	"github.com/scottylabs/cmuinsta/backend/middleware"
	"github.com/scottylabs/cmuinsta/backend/router"
)

func main() {
	db.Init()

	if err := middleware.Init(context.Background()); err != nil {
		log.Fatal("Failed to initialize OIDC provider:", err)
	}

	r := router.Setup()
	r.Run(":" + os.Getenv("BACKEND_PORT"))
}
