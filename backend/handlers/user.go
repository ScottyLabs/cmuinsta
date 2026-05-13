package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scottylabs/cmuinsta/backend/db"
	"github.com/scottylabs/cmuinsta/backend/models"
	"gorm.io/gorm"
)

// CreateUser creates an entry in the database for the user described by the
// provided gin.Context object. It also instantiates the directories to store
// user data. This covers the endpoint in RFC 2.1 `POST /me`.
// Authentication is handled by [middleware.Authenticate] before this occurs.
func CreateUser(c *gin.Context) {
	andrewID := c.GetString("andrewid")

	var body struct {
		InstagramID string `json:"instagram_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Major       string `json:"major" binding:"required"`
		Hometown    string `json:"hometown" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is malformed and missing components",
		})
		return
	}

	// Lookup by UUID
	ctx := c.Request.Context()
	instaPair, err := models.UUIDLookup(ctx, db.DB, body.InstagramID)

	// Handle the several errors we can encounter during translation
	if err != nil || instaPair == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "Instagram account not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "The server encountered an error translating your UUID"})
		return
	}

	user := models.User{
		AndrewID:  andrewID,
		Name:      body.Name,
		Major:     body.Major,
		Hometown:  body.Hometown,
		Instagram: instaPair.Username,
		Queued:    nil,
		Position:  -1,
		Posted:    nil,
	}

	// Instantiate directories
	store_path := os.Getenv("POSTS_STORE_DIR")
	// Create top-level user directory @ POSTS_STORE_DIR/[andrewid]
	cwd := fmt.Sprintf("%s/%s", store_path, andrewID)
	if err := os.MkdirAll(cwd, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize user storage"})
		return
	}
	// Create images directory @ POSTS_STORE_DIR/[andrewid]/images
	cwd = fmt.Sprintf("%s/%s/images", store_path, andrewID)
	if err := os.MkdirAll(cwd, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize user storage"})
		return
	}
	// Touch an empty description text file @ POSTS_STORE_DIR/[andrewid]/desc.txt
	captionPath := fmt.Sprintf("%s/%s/caption.txt", store_path, andrewID)
	if _, err := os.Create(captionPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize caption file"})
		return
	}

	// Database trasncation occurs (filesystem has passed)
	err = db.DB.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			err := instaPair.Consume(ctx, tx)
			if err != nil {
				return err
			}
			return tx.Create(&user).Error
		})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"andrewid":       user.AndrewID,
			"name":           user.Name,
			"major":          user.Major,
			"hometown":       user.Hometown,
			"instagram":      user.Instagram,
			"caption":        "",
			"image_count":    0,
			"queued_at":      nil,
			"queue_position": -1,
			"posted_at":      nil,
		})
}

// GetUser reqtrieves a user's information from the dataabse and filesystem
// based on the authenticated JWT. Returns it via Gin as an JSON object as
// specified by RFC 2.3 for the `GET /me` endpoint.
// Authentication is handled by [middleware.Authenticate] before this occurs.
func GetUser(c *gin.Context) {
	andrewID := c.GetString("andrewid")
	ctx := c.Request.Context()
	user, err := models.AndrewIDLookup(ctx, db.DB, andrewID)

	// Get user record
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "User not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	// Compute Queue Position
	pos := -1
	if user.Position != -1 {
		appPos, err := models.GetAppPos(ctx, db.DB)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to retrieve queue state"})
			return
		}
		pos = user.Position - appPos
	}

	// Description read
	store_path := os.Getenv("POSTS_STORE_DIR")
	// Read caption from filesystem, capped at 2200 characters
	caption := ""
	captionPath := fmt.Sprintf("%s/%s/caption.txt", store_path, andrewID)
	captionBytes, err := os.ReadFile(captionPath)
	if err == nil {
		runes := []rune(string(captionBytes))
		if len(runes) > 2200 {
			runes = runes[:2200]
		}
		caption = string(runes)
	}

	// Count of number of images
	imageCount := 0
	imagesDir := fmt.Sprintf("%s/%s/images", store_path, andrewID)
	entries, err := os.ReadDir(imagesDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				imageCount++
			}
		}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"andrewid":       user.AndrewID,
			"name":           user.Name,
			"major":          user.Major,
			"hometown":       user.Hometown,
			"instagram":      user.Instagram,
			"caption":        caption,
			"image_count":    imageCount,
			"queued_at":      user.Queued,
			"queue_position": pos,
			"posted_at":      user.Posted,
		})

}
