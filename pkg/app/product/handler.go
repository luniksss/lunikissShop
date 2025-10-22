package product

import (
	"encoding/json"
	"net/http"

	"lunikissShop/pkg/domain/product"
)

type Handler struct {
	service *SalesOutletService
}

func NewHandler(service *SalesOutletService) *Handler {
	return &Handler{
		service: service,
	}
}

// GetSalesOutlet возвращает список всех складов
// @Summary Получить список складов
// @Description Возвращает список всех активных складов
// @Tags outlet
// @Accept json
// @Produce json
// @Success 200 {object} SalesOutletResponse
// @Router /api/outlet/list [get]
func (h *Handler) GetSalesOutlet(w http.ResponseWriter, r *http.Request) {
	salesOutlets, err := h.service.GetAllSalesOutlet(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SalesOutletResponse{
		Success: true,
		Data:    salesOutlets,
		Count:   len(salesOutlets),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSalesOutletProducts возвращает все товары на складе
// @Summary Получить товары на складе
// @Description Возвращает список всех товаров на указанном складе
// @Tags warehouse
// @Accept json
// @Produce json
// @Param outlet_id path string true "ID склада"
// @Success 200 {object} ProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/outlet/{outlet_id}/product/list [get]
func (h *Handler) GetSalesOutletProducts(w http.ResponseWriter, r *http.Request) {
	salesOutletID := r.PathValue("outlet_id")
	if salesOutletID == "" {
		http.Error(w, "Sales Outlet ID is required", http.StatusBadRequest)
		return
	}

	var products []product.StockItem
	var err error
	products, err = h.service.GetAllSalesOutletProducts(r.Context(), salesOutletID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProductsResponse{
		Success: true,
		Data:    products,
		Count:   len(products),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
