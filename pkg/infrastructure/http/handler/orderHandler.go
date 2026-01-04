package handler

import (
	"encoding/json"
	"net/http"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/domain/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var orderInfo model.OrderRequestInfo
	if err := json.NewDecoder(r.Body).Decode(&orderInfo); err != nil {
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

	if err := h.orderService.DeleteOrderItem(ctx, orderItemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.PathValue("orderID")

	if err := h.orderService.DeleteOrder(ctx, orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
