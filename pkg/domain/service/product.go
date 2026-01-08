package service

import (
	"context"
	"errors"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/infrastructure/mysql/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo}
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	return ps.productRepo.GetAllProducts(ctx)
}

func (ps *ProductService) GetProduct(ctx context.Context, productID string) (model.Product, error) {
	return ps.productRepo.GetProductByID(ctx, productID)
}

func (ps *ProductService) AddProduct(ctx context.Context, product model.Product) error {
	existingProduct, err := ps.productRepo.GetProductByName(ctx, product.Name)

	if err == nil && existingProduct.ID != "" {
		return errors.New("product with this name already exists")
	}

	productId, err := ps.productRepo.AddProduct(ctx, &product)
	if err != nil {
		return err
	}

	product.Image.ProductID = productId
	err = ps.productRepo.AddProductImage(ctx, &product.Image)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, productInfo model.Product) error {
	product, err := ps.productRepo.GetProductByID(ctx, productInfo.ID)
	if err != nil || product.ID == "" {
		return errors.New("product not found")
	}
	if productInfo == product {
		return errors.New("product is already updated")
	}

	err = ps.productRepo.UpdateProduct(ctx, &productInfo)
	if err != nil {
		return err
	}

	if productInfo.Image != product.Image {
		err = ps.productRepo.UpdateProductImage(ctx, &productInfo.Image)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productID string) error {
	product, err := ps.productRepo.GetProductByID(ctx, productID)
	if err != nil || product.ID == "" {
		return errors.New("product not found")
	}

	err = ps.productRepo.DeleteProductImage(ctx, productID)
	if err != nil {
		return err
	}
	err = ps.productRepo.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}
