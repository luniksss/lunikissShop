package product

import (
	"context"
	"errors"

	"lunikissShop/pkg/domain/product"
)

type SalesOutletService struct {
	repo product.SalesOutletRepository
}

func NewSalesOutletService(repo product.SalesOutletRepository) *SalesOutletService {
	return &SalesOutletService{repo: repo}
}

// GetAllSalesOutlet возвращает все склады
func (s *SalesOutletService) GetAllSalesOutlet(ctx context.Context) ([]product.SalesOutlet, error) {
	return s.repo.GetAllSalesOutlet(ctx)
}

// GetSalesOutletByID возвращает информацию о складе
func (s *SalesOutletService) GetSalesOutletByID(ctx context.Context, salesOutletID string) (*product.SalesOutlet, error) {
	return s.repo.GetSalesOutletByID(ctx, salesOutletID)
}

// GetAllSalesOutletProducts возвращает все товары на складе (включая нулевые остатки)
func (s *SalesOutletService) GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]product.StockItem, error) {
	if salesOutletID == "" {
		return nil, errors.New("salesOutletID is empty")
	}

	return s.repo.GetAllSalesOutletProducts(ctx, salesOutletID)
}

// GetSalesOutletProduct возвращает информацию о конкретном товаре на складе конкретной точки
func (s *SalesOutletService) GetSalesOutletProduct(ctx context.Context, salesOutletID, productID string) ([]product.StockItem, error) {
	if productID == "" || salesOutletID == "" {
		return nil, errors.New("productID or salesOutletID is empty")
	}

	return s.repo.GetProductStock(ctx, salesOutletID, productID)
}

// AddStockItem добавляет элемент в таблицу salesOutlet
func (s *SalesOutletService) AddStockItem(ctx context.Context, salesOutletID string, stockItem *product.StockItem) error {
	if salesOutletID == "" {
		return errors.New("salesOutletID is empty")
	}

	return s.repo.AddStockItem(ctx, salesOutletID, stockItem)
}

// UpdateStock обновляет количество товара на складе (для продавца/админа)
func (s *SalesOutletService) UpdateStock(ctx context.Context, salesOutletID, productID string, quantity int) error {
	// Проверка прав доступа
	return s.repo.UpdateStockQuantity(ctx, salesOutletID, productID, quantity)
}

// DeleteStockItem удаляет элемент из таблицы salesOutlet
func (s *SalesOutletService) DeleteStockItem(ctx context.Context, salesOutletID, productID string) error {
	if salesOutletID == "" {
		return errors.New("salesOutletID is empty")
	}

	return s.repo.DeleteStockItem(ctx, salesOutletID, productID)
}
