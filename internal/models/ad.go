package models

import (
	"time"

	"gorm.io/gorm"
)

type Ad struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"size:1000;not null" json:"description"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	Price       float64        `gorm:"not null" json:"price"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
