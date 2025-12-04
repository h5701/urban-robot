package models

import "time"

type Order struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    int64       `json:"user_id" gorm:"index"`
	Status    string      `json:"status"` // e.g. Pending, Shipped, Cancelled
	Total     int64       `json:"total"`  // cents
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}
