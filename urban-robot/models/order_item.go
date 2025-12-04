package models

type OrderItem struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	OrderID   uint   `json:"order_id" gorm:"index"`
	ProductID uint   `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"` // cents at time of purchase
	Product   Product `gorm:"foreignKey:ProductID"` // optional preload
}
