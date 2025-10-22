package http

import (
	"net/http"

	"lunikissShop/pkg/app/product"
)

func SetupRoutes(warehouseHandler *product.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/outlet/list", warehouseHandler.GetSalesOutlet)
	mux.HandleFunc("GET /api/outlet/{outlet_id}/product/list", warehouseHandler.GetSalesOutletProducts)

	//mux.Handle("PUT /api/warehouses/{warehouse_id}/stock",
	//	authMiddleware.RequireAuth(http.HandlerFunc(warehouseHandler.UpdateStock)))

	return mux
}
