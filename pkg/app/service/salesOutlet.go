package service

import (
	"context"
	"lunikissShop/pkg/domain/model"
)

type SalesOutletService interface {
	GetAllSalesOutlet(ctx context.Context) ([]model.SalesOutlet, error)
	GetSalesOutlet(ctx context.Context, id string) (model.SalesOutlet, error)
	AddSalesOutlet(ctx context.Context, address string) error
	UpdateSalesOutlet(ctx context.Context, salesOutletID, address string) error
	DeleteSalesOutlet(ctx context.Context, salesOutletID string) error
	GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]model.StockItem, error)
	GetProductStock(ctx context.Context, salesOutletID, productID string) ([]model.StockItem, error)
	AddStockItem(ctx context.Context, stockItem *model.StockItem) error
	UpdateStockAmount(ctx context.Context, salesOutletID, productID string, amount, size int) error
	DeleteStockItem(ctx context.Context, salesOutletID, productID string) error
}
