package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username" binding:"required,min=3,max=30"`
	Password  string    `gorm:"_" json:"-" binding:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Listing struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Title       string    `gorm:"not null" json:"title" binding:"required,min=5,max=100"`
	Description string    `gorm:"not null" json:"description" binding:"required,min=10,max=500"`
	Price       float64   `gorm:"not null" json:"price" binding:"required,gt=0"`
	ImageURL    string    `gorm:"not null" json:"image_url" binding:"required,url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
