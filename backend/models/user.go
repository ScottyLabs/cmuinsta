package models

import "time"

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
