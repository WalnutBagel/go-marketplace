package models

import (
	"time"

	"gorm.io/gorm"
)

type Ad struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`        // заголовок, максимум 100 символов
	Description string         `gorm:"size:1000;not null" json:"description"` // описание, максимум 1000 символов
	ImageURL    string         `gorm:"size:255" json:"image_url"`             // URL картинки
	Price       float64        `gorm:"not null" json:"price"`                 // цена
	UserID      uint           `gorm:"not null" json:"user_id"`               // id автора
	User        User           `gorm:"foreignKey:UserID" json:"user"`         // связь с моделью User
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
