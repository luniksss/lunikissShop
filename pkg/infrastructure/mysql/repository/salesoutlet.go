package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"lunikissShop/pkg/domain/model"
)

type SalesOutletRepository struct {
	db *sql.DB
}

func NewSalesOutletRepository(db *sql.DB) *SalesOutletRepository {
	return &SalesOutletRepository{
		db: db,
	}
}

func (sor *SalesOutletRepository) GetAllSalesOutlet(ctx context.Context) ([]model.SalesOutlet, error) {
	query := `
        SELECT id, address
        FROM sales_outlet 
        ORDER BY address
    `

	return sor.querySalesOutlet(ctx, query)
}

func (sor *SalesOutletRepository) GetSalesOutletByID(ctx context.Context, id string) (model.SalesOutlet, error) {
	query := `
        SELECT id, address
        FROM sales_outlet 
        WHERE id = ?
    `

	var so model.SalesOutlet
	err := sor.db.QueryRowContext(ctx, query, id).Scan(
		&so.ID, &so.Address,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return so, errors.New("sales outlet not found")
	}
	if err != nil {
		return so, err
	}

	return so, nil
}

func (sor *SalesOutletRepository) GetSalesOutletByName(ctx context.Context, name string) (model.SalesOutlet, error) {
	query := `
        SELECT id, address
        FROM sales_outlet 
        WHERE address = ?
    `

	var so model.SalesOutlet
	err := sor.db.QueryRowContext(ctx, query, name).Scan(
		&so.ID, &so.Address,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return so, errors.New("sales outlet not found")
	}
	if err != nil {
		return so, err
	}

	return so, nil
}

func (sor *SalesOutletRepository) AddSalesOutlet(ctx context.Context, address string) error {
	query := `INSERT INTO sales_outlet(address) VALUE ?`
	_, err := sor.db.ExecContext(ctx, query, address)
	if err != nil {
		return err
	}
	return nil
}

func (sor *SalesOutletRepository) UpdateSalesOutlet(ctx context.Context, salesOutletID, address string) error {
	query := `UPDATE sales_outlet 
		SET address = ? 
		WHERE id = ?
	`
	_, err := sor.db.ExecContext(ctx, query, address, salesOutletID)
	if err != nil {
		return err
	}

	return nil
}

func (sor *SalesOutletRepository) DeleteSalesOutlet(ctx context.Context, salesOutletID string) error {
	query := `DELETE FROM sales_outlet WHERE id = ?`
	_, err := sor.db.ExecContext(ctx, query, salesOutletID)
	if err != nil {
		return fmt.Errorf("failed to delete sales outlet: %w", err)
	}

	return nil
}

func (sor *SalesOutletRepository) GetAllSalesOutletProducts(ctx context.Context, salesOutletID string) ([]model.StockItem, error) {
	query := `
        SELECT 
            p.id, p.name, p.description, p.price, 
            pi.image_path, pi.id,
            s.sales_outlet_id, s.size, s.amount
        FROM product p
        INNER JOIN product_stock s ON p.id = s.product_id
        JOIN product_image pi ON p.id = pi.product_id
		WHERE s.sales_outlet_id = ?
        ORDER BY p.name
    `

	return sor.queryProducts(ctx, query, salesOutletID)
}

func (sor *SalesOutletRepository) GetProductStock(ctx context.Context, salesOutletID, productID string) ([]model.StockItem, error) {
	query := `
 		SELECT 
			p.id, p.name, p.description, p.price,
			pi.image_path, pi.id,
			s.sales_outlet_id, s.size, s.amount
		FROM product p
 		INNER JOIN product_stock s ON p.id = s.product_id
 		JOIN product_image pi ON p.id = pi.product_id
 		WHERE pi.product_id = ?
 		AND s.product_id = ?
 		ORDER BY p.name
	`

	return sor.queryConcreteProducts(ctx, query, salesOutletID, productID)
}

func (sor *SalesOutletRepository) AddStockItem(ctx context.Context, stockItem *model.StockItem) error {
	query := `INSERT INTO product_stock (product_id, size, amount) VALUES (?, ?, ?)`
	_, err := sor.db.ExecContext(ctx, query, stockItem.Product.ID, stockItem.Size, stockItem.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (sor *SalesOutletRepository) UpdateStockAmount(ctx context.Context, salesOutletID, productID string, amount, size int) error {
	query := `UPDATE product_stock 
		SET amount = ? 
		WHERE sales_outlet_id = ? 
		AND product_id = ? 
	  	AND size = ?
	`
	_, err := sor.db.ExecContext(ctx, query, amount, salesOutletID, productID, size)
	if err != nil {
		return err
	}

	return nil
}

func (sor *SalesOutletRepository) DeleteStockItem(ctx context.Context, salesOutletID, productID string) error {
	query := `DELETE FROM product_stock WHERE sales_outlet_id = ? AND product_id = ?`
	_, err := sor.db.ExecContext(ctx, query, salesOutletID, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product from sales outlet stock: %w", err)
	}

	return nil
}

func (sor *SalesOutletRepository) querySalesOutlet(ctx context.Context, query string) ([]model.SalesOutlet, error) {
	rows, err := sor.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []model.SalesOutlet
	for rows.Next() {
		var so model.SalesOutlet
		err := rows.Scan(&so.ID, &so.Address)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, so)
	}

	return warehouses, nil
}

// TODO объединить с queryConcreteProducts
func (sor *SalesOutletRepository) queryProducts(ctx context.Context, query, warehouseID string) ([]model.StockItem, error) {
	rows, err := sor.db.QueryContext(ctx, query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.StockItem
	for rows.Next() {
		var ps model.StockItem
		err := rows.Scan(
			&ps.Product.ID, &ps.Product.Name, &ps.Product.Description, &ps.Product.Price,
			&ps.Product.Image.ImagePath, &ps.Product.Image.ID,
			&ps.SalesOutletID, &ps.Size, &ps.Amount,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, ps)
	}

	return products, nil
}

func (sor *SalesOutletRepository) queryConcreteProducts(ctx context.Context, query, warehouseID, productID string) ([]model.StockItem, error) {
	rows, err := sor.db.QueryContext(ctx, query, productID, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.StockItem
	for rows.Next() {
		var ps model.StockItem
		err := rows.Scan(
			&ps.Product.ID, &ps.Product.Name, &ps.Product.Description, &ps.Product.Price,
			&ps.Product.Image.ImagePath, &ps.Product.Image.ID,
			&ps.SalesOutletID, &ps.Size, &ps.Amount,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, ps)
	}

	return products, nil
}
