package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lunikissShop/pkg/domain/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	query := `
        SELECT 
            p.id, p.name, p.description, p.price, 
            pi.image_path, pi.id
        FROM product p
        LEFT JOIN product_image pi ON p.id = pi.product_id
        ORDER BY p.name
    `

	return r.queryAllProducts(ctx, query)
}

func (r *ProductRepository) GetProductByID(ctx context.Context, productID string) (model.Product, error) {
	query := `
 		SELECT 
			p.id, p.name, p.description, p.price,
			pi.image_path, pi.id
		FROM product p
 		JOIN product_image pi ON p.id = pi.product_id
 		WHERE pi.product_id = ?
	`

	return r.queryConcreteProduct(ctx, query, productID)
}

func (r *ProductRepository) GetProductByName(ctx context.Context, productName string) (model.Product, error) {
	query := `
 		SELECT 
			p.id, p.name, p.description, p.price,
			pi.image_path, pi.id
		FROM product p
 		JOIN product_image pi ON p.id = pi.product_id
 		WHERE p.name = ?
	`
	return r.queryConcreteProduct(ctx, query, productName)
}

func (r *ProductRepository) AddProduct(ctx context.Context, productInfo *model.Product) error {
	query := `INSERT INTO product (id, name, description, price) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, productInfo.ID, productInfo.Name, productInfo.Description, productInfo.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) AddProductImage(ctx context.Context, productImageInfo *model.Image) error {
	query := `INSERT INTO product_image (id, product_id, image_path) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, productImageInfo.ID, productImageInfo.ProductID, productImageInfo.ImagePath)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, newProductInfo *model.Product) error {
	query := `UPDATE product 
		SET name = ?, description = ?, price = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, newProductInfo.Name, newProductInfo.Description, newProductInfo.Price, newProductInfo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) UpdateProductImage(ctx context.Context, newProductInfo *model.Image) error {
	query := `UPDATE product_image 
		SET image_path = ?
		WHERE product_id = ?
		AND id = ?
	`
	_, err := r.db.ExecContext(ctx, query, newProductInfo.ImagePath, newProductInfo.ProductID, newProductInfo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	query := `DELETE FROM product WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *ProductRepository) DeleteProductImage(ctx context.Context, productID string) error {
	query := `DELETE FROM product_image WHERE product_id = ?`
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *ProductRepository) queryAllProducts(ctx context.Context, query string) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price,
			&p.Image.ImagePath, &p.Image.ID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) queryConcreteProduct(ctx context.Context, query, productParam string) (model.Product, error) {
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, productParam).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price,
		&p.Image.ImagePath, &p.Image.ID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return p, errors.New("product not found")
	}
	if err != nil {
		return p, err
	}

	return p, nil
}
