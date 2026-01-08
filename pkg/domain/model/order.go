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
	ProductID string `json:"ProductID"`
	Amount    int    `json:"Amount"`
	Price     int    `json:"Price"`
	Size      int    `json:"Size"`
}

type OrderRequestInfo struct {
	UserID        string                 `json:"UserID"`
	SalesOutletID string                 `json:"SalesOutletID"`
	OrderItems    []OrderItemRequestInfo `json:"OrderItems"`
}

type OrderItemResponseInfo struct {
	ID           string `json:"id"`
	OrderID      string `json:"order_id"`
	ProductID    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	Amount       int    `json:"amount"`
	Price        int    `json:"price"`
	Size         int    `json:"size"`
}

type OrderResponseInfo struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	SalesOutletID      string    `json:"sales_outlet_id"`
	SalesOutletAddress string    `json:"address"`
	CreatedAt          time.Time `json:"created_at"`
	StatusName         string    `json:"status_name"`
}

type OrderRepository interface {
	ListOrders(ctx context.Context) ([]Order, error)
	ListUserOrders(ctx context.Context, userID string) ([]OrderResponseInfo, error)
	ListOrdersBySalesOutlet(ctx context.Context, salesOutletID string) ([]Order, error)
	GetOrder(ctx context.Context, orderID string) (Order, error)
	GetOrderByID(ctx context.Context, orderID string) ([]OrderItemResponseInfo, error)
	CreateOrder(ctx context.Context, orderInfo OrderRequestInfo) error
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error
	DeleteOrderItem(ctx context.Context, orderItemID string) error
	DeleteOrder(ctx context.Context, orderID string) error
}
