package service

import (
	"errors"
	"fmt"

	"futuremarket/models"
	"futuremarket/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderService holds repositories needed for order logic.
type OrderService struct {
	OrderRepo   repository.OrderRepo
	CartRepo    repository.CartRepo
	ProductRepo repository.ProductRepo
}

// Checkout converts authenticated user's cart into an order using a DB transaction.
func (s OrderService) Checkout(userID int64) error {
	// We need a DB handle to start a transaction. Use one of the repos' DB.
	db := s.OrderRepo.DB
	if db == nil {
		return errors.New("database not configured for OrderRepo")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 1) Load the user's cart
		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cart not found for user %d", userID)
			}
			return err
		}

		// 2) Load cart items (and ensure we have rows)
		var items []models.CartItem
		if err := tx.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("cart is empty")
		}

		// We'll accumulate the total price (in cents)
		var total int64 = 0

		// 3) Re-check stock with row locking and deduct stock
		for _, ci := range items {
			// Lock the product row FOR UPDATE
			var prod models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ci.ProductID).First(&prod).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("product %d not found", ci.ProductID)
				}
				return err
			}

			// Lock and check stock from Stock table
			var stock models.Stock
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ?", ci.ProductID).First(&stock).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("stock record not found for product %d", ci.ProductID)
				}
				return err
			}

			// Check stock availability
			if stock.Quantity < int(ci.Quantity) {
				return fmt.Errorf("insufficient stock for product %d: available=%d requested=%d", prod.ID, stock.Quantity, ci.Quantity)
			}

			// Deduct stock
			stock.Quantity = stock.Quantity - int(ci.Quantity)
			if err := tx.Model(&models.Stock{}).Where("id = ?", stock.ID).Update("quantity", stock.Quantity).Error; err != nil {
				return err
			}

			// Add to total - use PriceCents if available, otherwise try Price field
			price := prod.PriceCents
			if price == 0 {
				// Try to get price from preloaded product or query
				if ci.Product.PriceCents > 0 {
					price = ci.Product.PriceCents
				} else {
					var p models.Product
					if err := tx.Where("id = ?", ci.ProductID).First(&p).Error; err != nil {
						return err
					}
					price = p.PriceCents
				}
			}
			total += price * ci.Quantity
		}

		// 4) Create order
		order := models.Order{
			UserID: userID,
			Status: "Pending",
			Total:  total,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 5) Move cart_items -> order_items
		orderItems := make([]models.OrderItem, 0, len(items))
		for _, ci := range items {
			// Get product price - use PriceCents field
			var unitPrice int64
			if ci.Product.PriceCents > 0 {
				unitPrice = ci.Product.PriceCents
			} else {
				var p models.Product
				if err := tx.Where("id = ?", ci.ProductID).First(&p).Error; err != nil {
					return err
				}
				unitPrice = p.PriceCents
			}

			oi := models.OrderItem{
				OrderID:   order.ID,
				ProductID: ci.ProductID,
				Quantity:  ci.Quantity,
				UnitPrice: unitPrice,
			}
			orderItems = append(orderItems, oi)
		}
		if len(orderItems) > 0 {
			if err := tx.Create(&orderItems).Error; err != nil {
				return err
			}
		}

		// 6) Clear the cart (delete cart items)
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// Optionally: delete the cart record itself or keep it (we'll keep)
		// Commit: returning nil signals commit
		return nil
	})
}

// ListOrders returns user's past orders (simple listing).
func (s OrderService) ListOrders(userID int64) ([]models.Order, error) {
	db := s.OrderRepo.DB
	if db == nil {
		return nil, errors.New("database not configured for OrderRepo")
	}
	var orders []models.Order
	if err := db.Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
