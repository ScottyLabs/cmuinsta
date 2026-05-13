package models

import (
	"context"

	"gorm.io/gorm"
)

type AppState struct {
	ID            uint `gorm:"primaryKey"`
	QueuePosition int
	QueueSize     int
}

// GetAppPos gets the current position of the queue
func GetAppPos(ctx context.Context, db *gorm.DB) (int, error) {
	var data AppState

	result := db.WithContext(ctx).First(&data, 1)

	if result.Error != nil {
		return -1, result.Error
	}

	return data.QueuePosition, nil
}
