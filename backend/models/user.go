package models

import (
	"time"
)

type Status string

const (
	StatusName     Status = "NAME"
	StatusMajor    Status = "MAJOR"
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

// Store defines the database operations the bot and handlers need.
// The GORM implementation lives in db/store.go.
type Store interface {
	// GetByIGSID looks up a student by their Instagram ID.
	// Returns nil, nil if not found — caller handles the new-user case.
	GetByIGSID(igsid string) (*User, error)

	// Save persists changes to a student row.
	Save(student *User) error

	// CreateUnverified inserts a new student row with only IGSID set.
	// Called when a user first DMs but hasn't linked their Andrew ID yet.
	CreateUnverified(igsid string) (*User, error)
}
