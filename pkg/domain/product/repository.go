package product

import "context"

type SalesOutletRepository interface {
	GetAllSalesOutlet(ctx context.Context) ([]SalesOutlet, error)
	GetSalesOutletByID(ctx context.Context, salesOutletID string) (*SalesOutlet, error)

	GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]StockItem, error)
	GetProductStock(ctx context.Context, salesOutletID, productID string) ([]StockItem, error)
	AddStockItem(ctx context.Context, salesOutletID string, stockItem *StockItem) error
	UpdateStockQuantity(ctx context.Context, salesOutletID, productID string, quantity int) error
	DeleteStockItem(ctx context.Context, salesOutletID, productID string) error
}
