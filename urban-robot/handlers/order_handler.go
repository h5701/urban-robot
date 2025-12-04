// this file is assigned to Hina,

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
