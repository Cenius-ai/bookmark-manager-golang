package models

import (
	"time"

	"gorm.io/gorm"
)

type Bookmark struct {
	ID        uint           `gorm:"primaryKey"`
	URL       string         `gorm:"not null"`
	Title     string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Tags      []Tag          `gorm:"many2many:bookmark_tags;"`
}

type Tag struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Bookmarks []Bookmark     `gorm:"many2many:bookmark_tags;"`
}
