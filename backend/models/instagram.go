package models

type InstagramIdentifierPair struct {
	Username string `gorm:"primaryKey"`
	UUID     string `gorm:"unique"`
	Consumed bool   `gorm:"default:false"`
}
