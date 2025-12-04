package models

import "gorm.io/gorm"

// Product is an item that can be listed, searched and bought.
type Product struct {
	gorm.Model
	Name          string  `json:"name" gorm:"size:255"`
	SKU           string  `json:"sku" gorm:"size:100;uniqueIndex"`
	Description   string  `json:"description" gorm:"type:text"`
	Category      string  `json:"category" gorm:"size:100"`
	PriceCents    int64   `json:"price_cents"` // store price in cents to avoid float issues
	Stock         int64   `json:"stock"`
	ImageURL      string  `json:"image_url" gorm:"size:500"`
	AverageRating float32 `json:"average_rating"`
	ReviewCount   int64   `json:"review_count"`
}


