package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/domain/service"
	"lunikissShop/pkg/middleware"
)

type SalesOutletHandler struct {
	salesOutletService *service.SalesOutletService
}

func NewSalesOutletHandler(salesOutletService *service.SalesOutletService) *SalesOutletHandler {
	return &SalesOutletHandler{salesOutletService: salesOutletService}
}

func (h *SalesOutletHandler) GetSalesOutlet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.salesOutletService.GetAllSalesOutlet(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *SalesOutletHandler) GetSalesOutletByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("id")

	product, err := h.salesOutletService.GetSalesOutlet(ctx, salesOutletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *SalesOutletHandler) AddSalesOutlet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var address string
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.salesOutletService.AddSalesOutlet(ctx, address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesOutletHandler) UpdateSalesOutlet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("outletId")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var address string
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.salesOutletService.UpdateSalesOutlet(ctx, salesOutletID, address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesOutletHandler) DeleteSalesOutlet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("outletId")

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := h.salesOutletService.DeleteSalesOutlet(ctx, salesOutletID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesOutletHandler) GetSalesOutletProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("id")

	products, err := h.salesOutletService.GetAllSalesOutletProducts(ctx, salesOutletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *SalesOutletHandler) GetSalesOutletProductsByProductID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("outletId")
	productID := r.PathValue("productId")

	product, err := h.salesOutletService.GetProductStock(ctx, salesOutletID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *SalesOutletHandler) AddStockItem(w http.ResponseWriter, r *http.Request) {
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

	var stockItem model.StockItem
	if err := json.NewDecoder(r.Body).Decode(&stockItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.salesOutletService.AddStockItem(ctx, &stockItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesOutletHandler) UpdateStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("outletId")
	productID := r.PathValue("productId")
	amount, _ := strconv.Atoi(r.PathValue("amount"))
	size, _ := strconv.Atoi(r.PathValue("size"))

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin && model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := h.salesOutletService.UpdateStockAmount(ctx, salesOutletID, productID, amount, size); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesOutletHandler) DeleteStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	salesOutletID := r.PathValue("outletId")
	productID := r.PathValue("productId")
	size, _ := strconv.Atoi(r.PathValue("size"))

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if model.Role(user.Role) != model.RoleAdmin && model.Role(user.Role) != model.RoleSeller {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if err := h.salesOutletService.DeleteStockItem(ctx, salesOutletID, productID, size); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
