package repository

import (
	"futuremarket/models"
	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

func (r OrderRepo) ListOrders(userID int64) ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
