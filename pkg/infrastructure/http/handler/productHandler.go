package handler

import (
	"encoding/json"
	"net/http"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/domain/service"
	"lunikissShop/pkg/middleware"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.PathValue("id")

	product, err := h.productService.GetProduct(ctx, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
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

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.productService.AddProduct(ctx, product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.productService.UpdateProduct(ctx, product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	productID := r.PathValue("id")
	if err := h.productService.DeleteProduct(ctx, productID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
