package models

import "time"

type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []CartItem `gorm:"foreignKey:CartID"`
}

