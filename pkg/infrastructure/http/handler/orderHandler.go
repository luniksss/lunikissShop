package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/domain/service"
	"lunikissShop/pkg/middleware"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin && model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	orders, err := h.orderService.ListAllOrders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) ListAllUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.PathValue("userID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if userID != user.ID &&
		model.Role(user.Role) != model.RoleAdmin &&
		model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	orders, err := h.orderService.ListAllUserOrders(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) ListOrdersBySalesOutlet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("salesOutletID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin && model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	orders, err := h.orderService.ListOrdersBySalesOutlet(ctx, salesOutletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.PathValue("orderID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !model.Role(user.Role).HasPermission(model.RoleUser) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	orderItems, err := h.orderService.GetOrderInfo(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderItems)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !model.Role(user.Role).HasPermission(model.RoleUser) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var orderInfo model.OrderRequestInfo
	if err := json.Unmarshal(bodyBytes, &orderInfo); err != nil {
		var raw map[string]interface{}
		if err2 := json.Unmarshal(bodyBytes, &raw); err2 != nil {
			log.Printf("Cannot parse as generic map: %v", err2)
		} else {
			for k, v := range raw {
				log.Printf("Key: %s, Type: %T, Value: %v", k, v, v)
			}
		}

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.orderService.CreateOrder(ctx, orderInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.PathValue("orderID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin && model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.orderService.UpdateOrderStatus(ctx, orderID, request.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderItemID := r.PathValue("orderItemID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !model.Role(user.Role).HasPermission(model.RoleUser) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := h.orderService.DeleteOrderItem(ctx, orderItemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.PathValue("orderID")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !model.Role(user.Role).HasPermission(model.RoleUser) {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := h.orderService.DeleteOrder(ctx, orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
