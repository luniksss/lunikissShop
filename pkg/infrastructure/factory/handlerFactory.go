package factory

import "lunikissShop/pkg/infrastructure/http/handler"

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
