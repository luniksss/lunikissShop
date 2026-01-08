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
	UserService        service.UserService
	AuthService        service.AuthService

	ProductHandler     *handler.ProductHandler
	SalesOutletHandler *handler.SalesOutletHandler
	OrderHandler       *handler.OrderHandler
	UserHandler        *handler.UserHandler
	AuthHandler        *handler.AuthHandler
}

func NewApp(db *sql.DB) *App {
	repoFactory := factory.NewRepositoryFactory(db)
	serviceFactory := factory.NewServiceFactory(*repoFactory)
	handlerFactory := factory.NewHandlerFactory(serviceFactory)

	productService := *serviceFactory.NewProductService()
	salesOutletService := *serviceFactory.NewSalesOutletService()
	orderService := *serviceFactory.NewOrderService(salesOutletService)
	userService := *serviceFactory.NewUserService()
	authService := *serviceFactory.NewAuthService(userService, "")

	productHandler := handlerFactory.ProductHandlers()
	salesOutletHandler := handlerFactory.SalesOutletHandlers()
	orderHandler := handlerFactory.OrderHandlers(salesOutletService)
	userHandler := handlerFactory.UserHandlers()
	authHandler := handlerFactory.AuthHandlers(userService, "")

	app := &App{
		RepoFactory:    repoFactory,
		ServiceFactory: serviceFactory,
		HandlerFactory: handlerFactory,

		ProductService:     productService,
		SalesOutletService: salesOutletService,
		OrderService:       orderService,
		UserService:        userService,
		AuthService:        authService,

		ProductHandler:     productHandler,
		SalesOutletHandler: salesOutletHandler,
		OrderHandler:       orderHandler,
		UserHandler:        userHandler,
		AuthHandler:        authHandler,
	}

	return app
}
