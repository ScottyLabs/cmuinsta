package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	Statusname     Status = "PENDING"
	StatusMajor    Status = "APPROVED"
	StatusHometown Status = "HOMETOWN"
	StatusPhoto    Status = "PHOTO"
	StatusCaption  Status = "CAPTION"
	StatusComplete Status = "COMPLETE"
)

type User struct {
	// User Data needed for the post
	AndrewID string `gorm:"primaryKey"`
	IGSID    string `gorm:"not null;unique"`
	Username string `gorm:"not null;unique"`
	Name     string
	Major    string
	Hometown string

	// Data needed for internal state maintenance
	State    Status `gorm:"not null"`
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

// IGSIDLookup looks up a user in the database based on their IGSID and
// returns the corresponding User struct.
func IGSIDLookup(ctx context.Context, db *gorm.DB, igsid string) (*User, error) {
	var user User

	result := db.WithContext(ctx).First(&user, "igsid = ?", igsid)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
