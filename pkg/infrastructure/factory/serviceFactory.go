package factory

import (
	"lunikissShop/pkg/domain/service"
)

type ServiceFactory struct {
	repoFactory RepositoryFactory
}

func NewServiceFactory(repoFactory RepositoryFactory) *ServiceFactory {
	return &ServiceFactory{repoFactory: repoFactory}
}

func (f *ServiceFactory) NewProductService() *service.ProductService {
	return service.NewProductService(*f.repoFactory.NewProductRepository())
}

func (f *ServiceFactory) NewSalesOutletService() *service.SalesOutletService {
	return service.NewSalesOutletService(*f.repoFactory.NewSalesOutletRepository(), *f.repoFactory.NewProductRepository())
}

func (f *ServiceFactory) NewOrderService(salesOutletService service.SalesOutletService) *service.OrderService {
	return service.NewOrderService(*f.repoFactory.NewOrderRepository(), salesOutletService)
}

func (f *ServiceFactory) NewUserService() *service.UserService {
	return service.NewUserService(*f.repoFactory.NewUserRepository())
}
