package service

import (
	"context"
	"lunikissShop/pkg/domain/model"
)

type OrderService interface {
	ListAllOrders(ctx context.Context) ([]model.Order, error)
	ListAllUserOrders(ctx context.Context, userID string) ([]model.Order, error)
	ListOrdersBySalesOutlet(ctx context.Context, salesOutletID string) ([]model.Order, error)
	GetOrderInfo(ctx context.Context, orderID string) ([]model.OrderItem, error)
	CreateOrder(ctx context.Context, orderInfo model.OrderRequestInfo) error
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error
	DeleteOrderItem(ctx context.Context, orderItemID string) error
	DeleteOrder(ctx context.Context, orderID string) error
}
