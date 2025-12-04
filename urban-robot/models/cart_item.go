package models

type CartItem struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	CartID    uint   `json:"cart_id" gorm:"index"`
	ProductID uint   `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID"` // preloadable
}
