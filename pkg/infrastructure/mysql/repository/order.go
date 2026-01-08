package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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
	query := "SELECT id, user_id, sales_outlet_id, created_at, status_name FROM `order`	ORDER BY id"
	rows, err := or.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var uo model.Order
		var createdAtBytes []byte
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.SalesOutletID, &createdAtBytes, &uo.StatusName)
		if err != nil {
			return nil, err
		}

		createdAtStr := string(createdAtBytes)
		uo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Failed to parse created_at: %v", err)
		}
		orders = append(orders, uo)
	}

	return orders, nil
}

func (or OrderRepository) ListUserOrders(ctx context.Context, userID string) ([]model.OrderResponseInfo, error) {
	query := "SELECT o.id, o.user_id, o.sales_outlet_id, o.created_at, o.status_name, so.address as address FROM `order` o LEFT JOIN sales_outlet so ON o.sales_outlet_id = so.id WHERE o.user_id = ? ORDER BY o.created_at DESC"

	rows, err := or.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderResponseInfo
	for rows.Next() {
		var uo model.OrderResponseInfo
		var createdAtBytes []byte
		err := rows.Scan(
			&uo.ID,
			&uo.UserID,
			&uo.SalesOutletID,
			&createdAtBytes,
			&uo.StatusName,
			&uo.SalesOutletAddress,
		)
		if err != nil {
			return nil, err
		}

		createdAtStr := string(createdAtBytes)
		uo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Failed to parse created_at: %v", err)
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
	query := "SELECT id, user_id, sales_outlet_id, created_at, status_name FROM `order`	WHERE id = ?"
	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var oi model.Order
	var createdAtBytes []byte
	for rows.Next() {
		err := rows.Scan(&oi.ID, &oi.UserID, &oi.SalesOutletID, &createdAtBytes, &oi.StatusName)
		if err != nil {
			return model.Order{}, err
		}

		createdAtStr := string(createdAtBytes)
		oi.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Failed to parse created_at: %v", err)
		}
	}

	return oi, nil
}

func (or OrderRepository) GetOrderByID(ctx context.Context, orderID string) ([]model.OrderItemResponseInfo, error) {
	query := `
    	SELECT 
            oi.id, 
            oi.order_id, 
            oi.product_id, 
            p.name as product_name,
            pi.image_path as product_image,
            oi.amount, 
            oi.price, 
            oi.size
        FROM order_item oi
        LEFT JOIN product p ON oi.product_id = p.id
        LEFT JOIN product_image pi ON pi.product_id = p.id
        WHERE oi.order_id = ?
        ORDER BY oi.id
	`
	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []model.OrderItemResponseInfo
	for rows.Next() {
		var oi model.OrderItemResponseInfo
		err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductName, &oi.ProductImage, &oi.Amount, &oi.Price, &oi.Size)
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

	query := "INSERT INTO `order` (user_id, sales_outlet_id, status_name) VALUES (?, ?, ?)"
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
		updateStockQuery := `
            UPDATE product_stock 
            SET amount = amount - ? 
            WHERE product_id = ? AND sales_outlet_id = ? AND size = ?
        `

		_, err := tx.ExecContext(ctx, updateStockQuery,
			orderItem.Amount,
			orderItem.ProductID,
			orderInfo.SalesOutletID,
			orderItem.Size,
		)

		if err != nil {
			return fmt.Errorf("не удалось уменьшить остатки: %w", err)
		}

		subQuery := `
            INSERT INTO order_item (order_id, product_id, amount, price, size)
            VALUES (?, ?, ?, ?, ?)
        `

		_, err = tx.ExecContext(ctx, subQuery,
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
	query := "UPDATE `order` SET status_name = ? WHERE id = ?"
	_, err := or.db.ExecContext(ctx, query, newStatus, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (or OrderRepository) DeleteOrderItem(ctx context.Context, orderItemID string) error {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := "SELECT o.sales_outlet_id, oi.product_id, oi.size, oi.amount, oi.order_id FROM order_item oi JOIN `order` o ON oi.order_id = o.id WHERE oi.id = ?"
	var (
		salesOutletID string
		productID     string
		size          int
		amount        int
		orderID       string
	)

	err = tx.QueryRowContext(ctx, query, orderItemID).Scan(
		&salesOutletID,
		&productID,
		&size,
		&amount,
		&orderID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order item not found: %w", err)
		}
		return fmt.Errorf("failed to query order item: %w", err)
	}

	updateStockQuery := `
        UPDATE product_stock 
        SET amount = amount + ? 
        WHERE product_id = ? 
          AND size = ?
          AND sales_outlet_id = ?
    `

	result, err := tx.ExecContext(ctx, updateStockQuery,
		amount,
		productID,
		size,
		salesOutletID,
	)

	if err != nil {
		return fmt.Errorf("failed to restore stock for product %s: %w", productID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Warning: failed to get rows affected for product %s: %v", productID, err)
	} else if rowsAffected == 0 {
		log.Printf("Warning: no stock record found for product %s, size %d, outlet %s",
			productID, size, salesOutletID)
	}

	deleteQuery := `DELETE FROM order_item WHERE id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (or OrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := "SELECT o.sales_outlet_id, oi.product_id, oi.size, oi.amount FROM `order` o JOIN order_item oi ON o.id = oi.order_id WHERE o.id = ?"
	rows, err := tx.QueryContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	type stockUpdate struct {
		salesOutletID string
		productID     string
		size          int
		amount        int
	}

	var updates []stockUpdate
	for rows.Next() {
		var update stockUpdate
		err := rows.Scan(&update.salesOutletID, &update.productID, &update.size, &update.amount)
		if err != nil {
			return fmt.Errorf("failed to scan order item: %w", err)
		}
		updates = append(updates, update)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating order items: %w", err)
	}

	updateStockQuery := `
        UPDATE product_stock 
        SET amount = amount + ? 
        WHERE product_id = ? 
          AND size = ?
          AND sales_outlet_id = ?
    `

	for _, update := range updates {
		result, err := tx.ExecContext(ctx, updateStockQuery,
			update.amount,
			update.productID,
			update.size,
			update.salesOutletID,
		)

		if err != nil {
			return fmt.Errorf("failed to restore stock for product %s: %w", update.productID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Warning: failed to get rows affected for product %s: %v", update.productID, err)
		} else if rowsAffected == 0 {
			log.Printf("Warning: no stock record found for product %s, size %d, outlet %s",
				update.productID, update.size, update.salesOutletID)
		}
	}

	deleteOrderItemsQuery := `DELETE FROM order_item WHERE order_id = ?`
	_, err = tx.ExecContext(ctx, deleteOrderItemsQuery, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order items: %w", err)
	}

	deleteOrderQuery := "DELETE FROM `order` WHERE id = ?"
	_, err = tx.ExecContext(ctx, deleteOrderQuery, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
