package repository

import (
	"context"
	"database/sql"
	"errors"

	"lunikissShop/pkg/domain/product"
)

type SalesOutletRepository struct {
	db *sql.DB
}

func NewSalesOutletRepository(db *sql.DB) *SalesOutletRepository {
	return &SalesOutletRepository{db: db}
}

func (r *SalesOutletRepository) GetAllSalesOutlet(ctx context.Context) ([]product.SalesOutlet, error) {
	query := `
        SELECT id, address
        FROM sales_outlet 
        ORDER BY name
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []product.SalesOutlet
	for rows.Next() {
		var so product.SalesOutlet
		err := rows.Scan(&so.ID, &so.Address)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, so)
	}

	return warehouses, nil
}

func (r *SalesOutletRepository) GetSalesOutletByID(ctx context.Context, id string) (*product.SalesOutlet, error) {
	query := `
        SELECT id, address
        FROM sales_outlet 
        WHERE id = ?
    `

	var so product.SalesOutlet
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&so.ID, &so.Address,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("sales outlet not found")
	}
	if err != nil {
		return nil, err
	}

	return &so, nil
}

func (r *SalesOutletRepository) GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]product.StockItem, error) {
	query := `
        SELECT 
            p.id, p.name, p.description, p.price,
            s.sale_outlet_id, s.size, s.amount
        FROM product p
        INNER JOIN product_stock s ON p.id = s.product_id
        WHERE s.sale_outlet_id = ?
        ORDER BY p.name
    `

	return r.queryProducts(ctx, query, salesOutletID)
}

func (r *SalesOutletRepository) queryProducts(ctx context.Context, query, warehouseID string) ([]product.StockItem, error) {
	rows, err := r.db.QueryContext(ctx, query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.StockItem
	for rows.Next() {
		var ps product.StockItem
		err := rows.Scan(
			&ps.Product.ID, &ps.Product.Name, &ps.Product.Description, &ps.Product.Price,
			&ps.SalesOutletID, &ps.Size, &ps.Amount,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, ps)
	}

	return products, nil
}
