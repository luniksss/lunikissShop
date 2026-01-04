package app

import (
	"database/sql"
	"lunikissShop/pkg/domain/service"
	"lunikissShop/pkg/infrastructure/factory"
	"lunikissShop/pkg/infrastructure/http/handler"
)

type App struct {
	RepoFactory    *factory.RepositoryFactory
	ServiceFactory *factory.ServiceFactory
	HandlerFactory *factory.HandlerFactory

	ProductService     service.ProductService
	SalesOutletService service.SalesOutletService
	OrderService       service.OrderService

	ProductHandler     *handler.ProductHandler
	SalesOutletHandler *handler.SalesOutletHandler
	OrderHandler       *handler.OrderHandler
}

func NewApp(db *sql.DB) *App {
	repoFactory := factory.NewRepositoryFactory(db)
	serviceFactory := factory.NewServiceFactory(*repoFactory)
	handlerFactory := factory.NewHandlerFactory(serviceFactory)

	productService := *serviceFactory.NewProductService()
	salesOutletService := *serviceFactory.NewSalesOutletService()
	orderService := *serviceFactory.NewOrderService(salesOutletService)

	productHandler := handlerFactory.ProductHandlers()
	salesOutletHandler := handlerFactory.SalesOutletHandlers()
	orderHandler := handlerFactory.OrderHandlers(salesOutletService)

	app := &App{
		RepoFactory:    repoFactory,
		ServiceFactory: serviceFactory,
		HandlerFactory: handlerFactory,

		ProductService:     productService,
		SalesOutletService: salesOutletService,
		OrderService:       orderService,

		ProductHandler:     productHandler,
		SalesOutletHandler: salesOutletHandler,
		OrderHandler:       orderHandler,
	}

	return app
}
