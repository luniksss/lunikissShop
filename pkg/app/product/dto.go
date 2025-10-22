package product

import "lunikissShop/pkg/domain/product"

type ProductsResponse struct {
	Success bool                `json:"success"`
	Data    []product.StockItem `json:"data"`
	Count   int                 `json:"count"`
}

type SalesOutletResponse struct {
	Success bool                  `json:"success"`
	Data    []product.SalesOutlet `json:"data"`
	Count   int                   `json:"count"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
