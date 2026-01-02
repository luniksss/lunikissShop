package model

import "context"

type SalesOutlet struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type StockItem struct {
	SalesOutletID string `json:"sales_outlet_id"`
	Product       Product
	Size          int `json:"size"`
	Amount        int `json:"amount"`
}

type SalesOutletRepository interface {
	GetAllSalesOutlet(ctx context.Context) ([]SalesOutlet, error)
	GetSalesOutletByID(ctx context.Context, id string) (SalesOutlet, error)
	GetSalesOutletByName(ctx context.Context, name string) (SalesOutlet, error)
	AddSalesOutlet(ctx context.Context, address string) error
	UpdateSalesOutlet(ctx context.Context, salesOutletID, address string) error
	DeleteSalesOutlet(ctx context.Context, salesOutletID string) error
	GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]StockItem, error)
	GetProductStock(ctx context.Context, salesOutletID, productID string) ([]StockItem, error)
	AddStockItem(ctx context.Context, stockItem *StockItem) error
	UpdateStockAmount(ctx context.Context, salesOutletID, productID string, amount, size int) error
	DeleteStockItem(ctx context.Context, salesOutletID, productID string) error
}
