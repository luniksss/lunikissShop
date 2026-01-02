package model

import "context"

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       Image   `json:"image"`
}

type Image struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	ImagePath string `json:"image_path"`
}

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, productID string) (Product, error)
	GetProductByName(ctx context.Context, productName string) (Product, error)
	AddProduct(ctx context.Context, productInfo *Product) error
	AddProductImage(ctx context.Context, productImageInfo *Image) error
	UpdateProduct(ctx context.Context, newProductInfo *Product) error
	UpdateProductImage(ctx context.Context, newProductInfo *Image) error
	DeleteProduct(ctx context.Context, productID string) error
	DeleteProductImage(ctx context.Context, productID string) error
}
