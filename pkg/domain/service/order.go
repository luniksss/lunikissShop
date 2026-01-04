package service

import (
	"context"
	"errors"

	"lunikissShop/pkg/app/service"
	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/infrastructure/mysql/repository"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	// TODO adapter to user service
	salesOutletService service.SalesOutletService
}

func NewOrderService(orderRepo repository.OrderRepository, salesOutletService service.SalesOutletService) *OrderService {
	return &OrderService{orderRepo, salesOutletService}
}

func (os OrderService) ListAllOrders(ctx context.Context) ([]model.Order, error) {
	return os.orderRepo.ListOrders(ctx)
}

func (os OrderService) ListAllUserOrders(ctx context.Context, userID string) ([]model.Order, error) {
	// TODO check if user with this ID exists

	return os.orderRepo.ListUserOrders(ctx, userID)
}

func (os OrderService) ListOrdersBySalesOutlet(ctx context.Context, salesOutletID string) ([]model.Order, error) {
	salesOutletExists := os.checkIsSalesOutletExists(ctx, salesOutletID)
	if !salesOutletExists {
		return []model.Order{}, errors.New("salesOutlet does not exist")
	}

	return os.orderRepo.ListOrdersBySalesOutlet(ctx, salesOutletID)
}

func (os OrderService) GetOrderInfo(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	return os.orderRepo.GetOrderByID(ctx, orderID)
}

func (os OrderService) CreateOrder(ctx context.Context, orderInfo model.OrderRequestInfo) error {
	salesOutlet := os.checkIsSalesOutletExists(ctx, orderInfo.SalesOutletID)
	user := os.checkIsUserExist(ctx, orderInfo.UserID)

	if !salesOutlet || !user {
		return errors.New("sales outlet or user does not exist")
	}

	return os.orderRepo.CreateOrder(ctx, orderInfo)
}

func (os OrderService) UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error {
	order, err := os.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.StatusName == newStatus {
		return errors.New("order status is same " + newStatus)
	}

	return os.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus)
}

func (os OrderService) DeleteOrderItem(ctx context.Context, orderItemID string) error {
	return os.orderRepo.DeleteOrderItem(ctx, orderItemID)
}

func (os OrderService) DeleteOrder(ctx context.Context, orderID string) error {
	orderExists := os.checkIsOrderExist(ctx, orderID)
	if !orderExists {
		return errors.New("order does not exist")
	}

	return os.orderRepo.DeleteOrder(ctx, orderID)
}

func (os OrderService) checkIsSalesOutletExists(ctx context.Context, salesOutletID string) bool {
	_, err := os.salesOutletService.GetSalesOutlet(ctx, salesOutletID)
	if err != nil {
		return false
	}
	return true
}

func (os OrderService) checkIsOrderExist(ctx context.Context, orderID string) bool {
	_, err := os.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return false
	}
	return true
}

func (os OrderService) checkIsUserExist(ctx context.Context, userID string) bool {
	return true
}
