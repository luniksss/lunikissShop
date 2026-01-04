package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lunikissShop/pkg/domain/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (or OrderRepository) ListOrders(ctx context.Context) ([]model.Order, error) {
	query := `
    	SELECT id, user_id, sales_outlet_id, created_at, status_name
    	FROM "order" 
    	ORDER BY id
	`
	rows, err := or.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var uo model.Order
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.SalesOutletID, &uo.CreatedAt, &uo.StatusName)
		if err != nil {
			return nil, err
		}
		orders = append(orders, uo)
	}

	return orders, nil
}

func (or OrderRepository) ListUserOrders(ctx context.Context, userID string) ([]model.Order, error) {
	query := `
    	SELECT id, user_id, sales_outlet_id, created_at, status_name
    	FROM "order" 
    	WHERE user_id = ?
    	ORDER BY id
	`
	rows, err := or.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var uo model.Order
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.SalesOutletID, &uo.CreatedAt, &uo.StatusName)
		if err != nil {
			return nil, err
		}
		orders = append(orders, uo)
	}

	return orders, nil
}

func (or OrderRepository) ListOrdersBySalesOutlet(ctx context.Context, salesOutletID string) ([]model.Order, error) {
	query := `
    	SELECT id, user_id, sales_outlet_id, created_at, status_name
    	FROM "order" 
    	WHERE sales_outlet_id = ?
    	ORDER BY id
	`
	rows, err := or.db.QueryContext(ctx, query, salesOutletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var uo model.Order
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.SalesOutletID, &uo.CreatedAt, &uo.StatusName)
		if err != nil {
			return nil, err
		}
		orders = append(orders, uo)
	}

	return orders, nil
}

func (or OrderRepository) GetOrder(ctx context.Context, orderID string) (model.Order, error) {
	query := `
    	SELECT id, user_id, id, user_id, sales_outlet_id, created_at, status_name
    	FROM "order" 
    	WHERE id = ?
	`
	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var oi model.Order
	for rows.Next() {
		err := rows.Scan(&oi.ID, &oi.UserID, &oi.SalesOutletID, &oi.CreatedAt, &oi.StatusName)
		if err != nil {
			return model.Order{}, err
		}
	}

	return oi, nil
}

func (or OrderRepository) GetOrderByID(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	query := `
    	SELECT id, order_id, product_id, amount, price, size
    	FROM order_item 
    	WHERE order_id = ?
    	ORDER BY id
	`
	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []model.OrderItem
	for rows.Next() {
		var oi model.OrderItem
		err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Amount, &oi.Price, &oi.Size)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, oi)
	}

	return orderItems, nil
}

func (or OrderRepository) CreateOrder(ctx context.Context, orderInfo model.OrderRequestInfo) error {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO "order" (user_id, sales_outlet_id, status_name)
        VALUES (?, ?, ?)
    `

	result, err := tx.ExecContext(ctx, query,
		orderInfo.UserID,
		orderInfo.SalesOutletID,
		"ordered",
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get order ID: %w", err)
	}

	for _, orderItem := range orderInfo.OrderItems {
		subQuery := `
            INSERT INTO order_item (order_id, product_id, amount, price, size)
            VALUES (?, ?, ?, ?, ?)
        `

		_, err := tx.ExecContext(ctx, subQuery,
			orderID,
			orderItem.ProductID,
			orderItem.Amount,
			orderItem.Price,
			orderItem.Size,
		)

		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (or OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error {
	query := `UPDATE "order" 
		SET status_name = ?
		WHERE id = ?
	`
	_, err := or.db.ExecContext(ctx, query, newStatus, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (or OrderRepository) DeleteOrderItem(ctx context.Context, orderItemID string) error {
	query := `DELETE FROM order_item WHERE id = ?`
	_, err := or.db.ExecContext(ctx, query, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	return nil
}

func (or OrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	query := `DELETE FROM "order" WHERE id = ?`
	_, err := or.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	query = `DELETE FROM order_item WHERE order_id = ?`
	_, err = or.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}
	return nil
}
