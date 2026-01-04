package service

import (
	"context"
	"errors"
	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/infrastructure/mysql/repository"
)

type SalesOutletService struct {
	outletRepo  repository.SalesOutletRepository
	productRepo repository.ProductRepository
}

func NewSalesOutletService(outletRepo repository.SalesOutletRepository, productRepo repository.ProductRepository) *SalesOutletService {
	return &SalesOutletService{outletRepo, productRepo}
}

func (s SalesOutletService) GetAllSalesOutlet(ctx context.Context) ([]model.SalesOutlet, error) {
	return s.outletRepo.GetAllSalesOutlet(ctx)
}

func (s SalesOutletService) GetSalesOutlet(ctx context.Context, id string) (model.SalesOutlet, error) {
	return s.findOutlet(ctx, id)
}

func (s SalesOutletService) AddSalesOutlet(ctx context.Context, address string) error {
	salesOutlet, err := s.findOutletByName(ctx, address)
	if err != nil {
		return err
	}
	if salesOutlet.Address != "" {
		return errors.New("sales outlet already exists")
	}

	return s.outletRepo.AddSalesOutlet(ctx, address)
}

func (s SalesOutletService) UpdateSalesOutlet(ctx context.Context, salesOutletID, address string) error {
	salesOutlet, err := s.findOutlet(ctx, salesOutletID)
	if err != nil {
		return err
	}
	if salesOutlet.Address != "" {
		return errors.New("sales outlet does not exists")
	}

	return s.outletRepo.UpdateSalesOutlet(ctx, salesOutletID, address)
}

func (s SalesOutletService) DeleteSalesOutlet(ctx context.Context, salesOutletID string) error {
	salesOutlet, err := s.findOutlet(ctx, salesOutletID)
	if err != nil {
		return err
	}
	if salesOutlet.Address != "" {
		return errors.New("sales outlet does not exists")
	}

	return s.outletRepo.DeleteSalesOutlet(ctx, salesOutletID)
}

func (s SalesOutletService) GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]model.StockItem, error) {
	_, err := s.findOutlet(ctx, salesOutletID)
	if err != nil {
		return []model.StockItem{}, err
	}

	return s.outletRepo.GetAllSalesOutletProducts(ctx, salesOutletID)
}

func (s SalesOutletService) GetProductStock(ctx context.Context, salesOutletID, productID string) ([]model.StockItem, error) {
	return s.isStockItemExists(ctx, salesOutletID, productID)
}

func (s SalesOutletService) AddStockItem(ctx context.Context, stockItem *model.StockItem) error {
	stockProducts, err := s.isStockItemExists(ctx, stockItem.SalesOutletID, stockItem.Product.ID)
	if err != nil {
		return err
	}
	if stockProducts != nil {
		for _, product := range stockProducts {
			if product.Size == stockItem.Size {
				return errors.New("product stock already exists")
			}
		}
	}

	return s.outletRepo.AddStockItem(ctx, stockItem)
}

func (s SalesOutletService) UpdateStockAmount(ctx context.Context, salesOutletID, productID string, amount, size int) error {
	stockProducts, err := s.isStockItemExists(ctx, salesOutletID, productID)
	if err != nil {
		return err
	}
	if stockProducts == nil {
		return errors.New("product stock does not exists")
	}

	return s.outletRepo.UpdateStockAmount(ctx, salesOutletID, productID, amount, size)
}

func (s SalesOutletService) DeleteStockItem(ctx context.Context, salesOutletID, productID string) error {
	stockProducts, err := s.isStockItemExists(ctx, salesOutletID, productID)
	if err != nil {
		return err
	}
	if stockProducts == nil {
		return errors.New("product stock does not exists")
	}

	return s.outletRepo.DeleteStockItem(ctx, salesOutletID, productID)
}

func (s SalesOutletService) isStockItemExists(ctx context.Context, salesOutletID, productID string) ([]model.StockItem, error) {
	_, err := s.findOutlet(ctx, salesOutletID)
	if err != nil {
		return []model.StockItem{}, err
	}

	_, err = s.findProduct(ctx, productID)
	if err != nil {
		return []model.StockItem{}, err
	}

	stockProducts, err := s.GetProductStock(ctx, salesOutletID, productID)
	if err != nil {
		return []model.StockItem{}, err
	}
	return stockProducts, nil
}

func (s SalesOutletService) findOutlet(ctx context.Context, salesOutletID string) (model.SalesOutlet, error) {
	outlet, err := s.outletRepo.GetSalesOutletByID(ctx, salesOutletID)
	if err != nil {
		return model.SalesOutlet{}, err
	}
	if outlet.ID == "" {
		return model.SalesOutlet{}, errors.New("outlet not found")
	}
	return outlet, nil
}

func (s SalesOutletService) findOutletByName(ctx context.Context, address string) (model.SalesOutlet, error) {
	outlet, err := s.outletRepo.GetSalesOutletByName(ctx, address)
	if err != nil {
		return model.SalesOutlet{}, err
	}
	if outlet.ID == "" {
		return model.SalesOutlet{}, errors.New("outlet not found")
	}
	return outlet, nil
}

func (s SalesOutletService) findProduct(ctx context.Context, productID string) (model.Product, error) {
	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return model.Product{}, err
	}
	if product.ID == "" {
		return model.Product{}, errors.New("product not found")
	}
	return product, nil
}
