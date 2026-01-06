package factory

import (
	"lunikissShop/pkg/domain/service"
	"lunikissShop/pkg/infrastructure/http/handler"
)

type HandlerFactory struct {
	serviceFactory *ServiceFactory
}

func NewHandlerFactory(serviceFactory *ServiceFactory) *HandlerFactory {
	return &HandlerFactory{serviceFactory: serviceFactory}
}

func (f *HandlerFactory) ProductHandlers() *handler.ProductHandler {
	return handler.NewProductHandler(f.serviceFactory.NewProductService())
}

func (f *HandlerFactory) SalesOutletHandlers() *handler.SalesOutletHandler {
	return handler.NewSalesOutletHandler(f.serviceFactory.NewSalesOutletService())
}

func (f *HandlerFactory) OrderHandlers(salesOutletService service.SalesOutletService) *handler.OrderHandler {
	return handler.NewOrderHandler(f.serviceFactory.NewOrderService(salesOutletService))
}

func (f *HandlerFactory) UserHandlers() *handler.UserHandler {
	return handler.NewUserHandler(f.serviceFactory.NewUserService())
}

func (f *HandlerFactory) AuthHandlers(userService service.UserService, jwtSecret string) *handler.AuthHandler {
	return handler.NewAuthHandler(f.serviceFactory.NewAuthService(userService, jwtSecret))
}
