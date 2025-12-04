// this file is assigned to Jasmine and Hina,

// PURPOSE:
// - HTTP handlers for checkout and order history.
// - Responsible for turning the cart into an order using DB transactions.
//
// EPICS & USER STORIES:
// - Epic 4: Checkout & Orders
//   - User Story 4.1: Place Order    (POST /api/v1/checkout)
//   - User Story 4.2: Order History  (GET /api/v1/orders)
//
// ENDPOINTS (to be implemented here):
// - POST /api/v1/checkout
//   - Uses the authenticated user's cart.
//   - Calls a service that performs an ACID transaction:
//     1) Re-check stock (with row locking if needed).
//     2) Deduct stock from products/stock table.
//     3) Create new order record.
//     4) Move cart_items → order_items.
//     5) Clear the cart.
//   - If ANY step fails, transaction MUST roll back.
//
// - GET /api/v1/orders
//   - Returns list of past orders for the authenticated user.
//   - Includes status (Pending, Shipped, Cancelled).
//

// Ladies just remember this,  Handler just needs::
//   - Gets user_id from JWT
//   - Calls orderService.Checkout(userID)
//   - Returns success/error response.

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

package handlers

import (
	"encoding/json"
	"futuremarket/middleware"
	"futuremarket/service"
	"net/http"
)

// OrderHandler manages checkout + history (Epic 4).
type OrderHandler struct {
	Service service.OrderService
}

// Helper: get authenticated user ID from context
func getUserID(r *http.Request) (int64, error) {
	val := r.Context().Value(middleware.ContextUserID)
	if val == nil {
		return 0, http.ErrNoCookie // unauthorized
	}
	// Middleware stores as int, convert to int64
	userID, ok := val.(int)
	if !ok {
		return 0, http.ErrNoCookie
	}
	return int64(userID), nil
}

// POST /api/v1/checkout
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	// Delegate ACID logic to the service layer
	if err := h.Service.Checkout(userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "order successfully created",
	})
}

// GET /api/v1/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	orders, err := h.Service.ListOrders(userID)
	if err != nil {
		http.Error(w, "failed to load orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
