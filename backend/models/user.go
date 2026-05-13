package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	// User Data needed for the post
	AndrewID  string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Major     string `gorm:"not null"`
	Hometown  string `gorm:"not null"`
	Instagram string `gorm:"not null;unique"`

	// Data needed for internal state maintenance
	Queued   *time.Time
	Position int `gorm:"default:-1"`
	Posted   *time.Time
}

// AndrewIDLookup looks up a user in the database based on their AndrewID and
// returns the corresponding User struct.
func AndrewIDLookup(ctx context.Context, db *gorm.DB, andrewid string) (*User, error) {
	var user User

	result := db.WithContext(ctx).First(&user, "andrewid = ?", andrewid)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
