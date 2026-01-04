package model

import (
	"context"
	"time"
)

type Order struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	SalesOutletID string    `json:"sales_outlet_id"`
	CreatedAt     time.Time `json:"created_at"`
	StatusName    string    `json:"status_name"`
}

type OrderItem struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Amount    int    `json:"amount"`
	Price     int    `json:"price"`
	Size      int    `json:"size"`
}

type OrderItemRequestInfo struct {
	ProductID string
	Amount    int
	Price     int
	Size      int
}

type OrderRequestInfo struct {
	UserID        string
	SalesOutletID string
	OrderItems    []OrderItemRequestInfo
}

type OrderRepository interface {
	ListOrders(ctx context.Context) ([]Order, error)
	ListUserOrders(ctx context.Context, userID string) ([]Order, error)
	ListOrdersBySalesOutlet(ctx context.Context, salesOutletID string) ([]Order, error)
	GetOrder(ctx context.Context, orderID string) (Order, error)
	GetOrderByID(ctx context.Context, orderID string) ([]OrderItem, error)
	CreateOrder(ctx context.Context, orderInfo OrderRequestInfo) error
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error
	DeleteOrderItem(ctx context.Context, orderItemID string) error
	DeleteOrder(ctx context.Context, orderID string) error
}
