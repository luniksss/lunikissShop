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

	ProductHandler     *handler.ProductHandler
	SalesOutletHandler *handler.SalesOutletHandler
}

func NewApp(db *sql.DB) *App {
	repoFactory := factory.NewRepositoryFactory(db)
	serviceFactory := factory.NewServiceFactory(*repoFactory)
	handlerFactory := factory.NewHandlerFactory(serviceFactory)

	app := &App{
		RepoFactory:    repoFactory,
		ServiceFactory: serviceFactory,
		HandlerFactory: handlerFactory,

		ProductService:     *serviceFactory.NewProductService(),
		SalesOutletService: *serviceFactory.NewSalesOutletService(),

		ProductHandler:     handlerFactory.ProductHandlers(),
		SalesOutletHandler: handlerFactory.SalesOutletHandlers(),
	}

	return app
}
