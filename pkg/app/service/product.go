package service

import (
	"context"
	"lunikissShop/pkg/domain/model"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, productID string) (model.Product, error)
	AddProduct(ctx context.Context, product model.Product) error
	UpdateProduct(ctx context.Context, product model.Product) error
	DeleteProduct(ctx context.Context, productID string) error
}
