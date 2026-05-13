package models

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type InstagramIdPair struct {
	Username string     `gorm:"primaryKey"`
	UUID     string     `gorm:"unique"`
	Consumed *time.Time `gorm:"default:false"`
}

// IsConsumed returns whether or not the account is consumed
func (a *InstagramIdPair) IsConsumed() bool {
	return a.Consumed != nil
}

// Consume marks the account as claimed by setting Consumed to now.
// Returns an error if the account is already consumed.
func (a *InstagramIdPair) Consume(ctx context.Context, db *gorm.DB) error {
	if a.IsConsumed() {
		return errors.New("account already consumed")
	}

	now := time.Now().UTC()

	// Use Model() to specify the target and Updates() for the change
	result := db.WithContext(ctx).
		Model(a).
		Where("uuid = ? AND consumed IS NULL", a.UUID).
		Update("consumed", now)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("failed to consume: account already consumed or not found")
	}

	// Update the local instance
	a.Consumed = &now
	return nil
}

// UUIDLookup converts a string UUID into the InstagramIdPair object. Returns
// nil if not found.
func UUIDLookup(ctx context.Context, db *gorm.DB, uuid string) (*InstagramIdPair, error) {
	var account InstagramIdPair

	result := db.WithContext(ctx).First(&account, "uuid = ?", uuid)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}
