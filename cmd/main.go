package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"lunikissShop/pkg/app"
	"lunikissShop/pkg/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	dbPath := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS lunishop")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database created successfully!")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	newApp := app.NewApp(db)
	setupRoutes(newApp)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With", "Accept"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	handler := c.Handler(http.DefaultServeMux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func setupRoutes(app *app.App) {
	protectedRouter := http.NewServeMux()

	http.HandleFunc("GET /product/list", app.ProductHandler.GetProducts)
	http.HandleFunc("GET /product/{id}", app.ProductHandler.GetProductByID)
	protectedRouter.HandleFunc("POST /api/v1/product/add", app.ProductHandler.AddProduct)
	protectedRouter.HandleFunc("POST /api/v1/product/update/{id}", app.ProductHandler.UpdateProduct)
	protectedRouter.HandleFunc("DELETE /api/v1/product/delete/{id}", app.ProductHandler.DeleteProduct)

	http.HandleFunc("GET /outlet/list", app.SalesOutletHandler.GetSalesOutlet)
	http.HandleFunc("GET /outlet/{id}", app.SalesOutletHandler.GetSalesOutletByID)
	protectedRouter.HandleFunc("POST /api/v1/outlet/add", app.SalesOutletHandler.AddSalesOutlet)
	protectedRouter.HandleFunc("POST /api/v1/outlet/update/{outletId}", app.SalesOutletHandler.UpdateSalesOutlet)
	protectedRouter.HandleFunc("DELETE /api/v1/outlet/delete/{outletId}", app.SalesOutletHandler.DeleteSalesOutlet)

	http.HandleFunc("GET /products/outlet/{id}", app.SalesOutletHandler.GetSalesOutletProducts)
	http.HandleFunc("GET /product/outlet/{outletId}/{productId}", app.SalesOutletHandler.GetSalesOutletProductsByProductID)
	protectedRouter.HandleFunc("POST /api/v1/stock/add", app.SalesOutletHandler.AddStockItem)
	protectedRouter.HandleFunc("POST /api/v1/stock/update/{outletId}/{productId}/{amount}/{size}", app.SalesOutletHandler.UpdateStockItem)
	protectedRouter.HandleFunc("DELETE /api/v1/stock/delete/{outletId}/{productId}/{size}", app.SalesOutletHandler.DeleteStockItem)

	protectedRouter.HandleFunc("GET /api/v1/orders/list", app.OrderHandler.ListAllOrders)
	protectedRouter.HandleFunc("GET /api/v1/orders/{orderID}", app.OrderHandler.GetOrderInfo)
	protectedRouter.HandleFunc("GET /api/v1/users/{userID}/orders", app.OrderHandler.ListAllUserOrders)
	protectedRouter.HandleFunc("GET /api/v1/sales-outlets/{salesOutletID}/orders", app.OrderHandler.ListOrdersBySalesOutlet)
	protectedRouter.HandleFunc("POST /api/v1/order", app.OrderHandler.CreateOrder)
	protectedRouter.HandleFunc("PATCH /api/v1/order/{orderID}/status", app.OrderHandler.UpdateOrderStatus)
	protectedRouter.HandleFunc("DELETE /api/v1/order/{orderID}", app.OrderHandler.DeleteOrder)
	protectedRouter.HandleFunc("DELETE /api/v1/order-items/{orderItemID}", app.OrderHandler.DeleteOrderItem)

	protectedRouter.HandleFunc("GET /api/v1/users", app.UserHandler.ListAllUsers)
	protectedRouter.HandleFunc("GET /api/v1/users/{id}", app.UserHandler.GetUserByID)
	protectedRouter.HandleFunc("POST /api/v1/users/by-email", app.UserHandler.GetUserByEmail)

	http.HandleFunc("POST /api/v1/users", app.UserHandler.AddUser)
	protectedRouter.HandleFunc("PUT /api/v1/users", app.UserHandler.UpdateUser)
	protectedRouter.HandleFunc("PATCH /api/v1/users/{id}/role", app.UserHandler.UpdateUserRole)
	protectedRouter.HandleFunc("DELETE /api/v1/users/{id}", app.UserHandler.DeleteUser)

	http.HandleFunc("POST /api/v1/auth/register", app.AuthHandler.Register)
	http.HandleFunc("POST /api/v1/auth/login", app.AuthHandler.Login)
	http.HandleFunc("POST /api/v1/auth/refresh", app.AuthHandler.RefreshToken)
	http.HandleFunc("POST /api/v1/auth/logout", app.AuthHandler.Logout)
	http.HandleFunc("POST /api/v1/auth/change-password", app.AuthHandler.ChangePassword)

	http.Handle("/api/v1/", middleware.AuthMiddleware(&app.AuthService)(protectedRouter))
}
