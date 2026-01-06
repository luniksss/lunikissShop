package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

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

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes(app *app.App) {
	protectedRouter := http.NewServeMux()

	http.HandleFunc("GET /product/list", app.ProductHandler.GetProducts)
	http.HandleFunc("GET /product/{id}", app.ProductHandler.GetProductByID)
	protectedRouter.HandleFunc("POST /product/add", app.ProductHandler.AddProduct)
	protectedRouter.HandleFunc("POST /product/update/{id}", app.ProductHandler.UpdateProduct)
	protectedRouter.HandleFunc("DELETE /product/delete/{id}", app.ProductHandler.DeleteProduct)

	http.HandleFunc("GET /outlet/list", app.SalesOutletHandler.GetSalesOutlet)
	http.HandleFunc("GET /outlet/{id}", app.SalesOutletHandler.GetSalesOutletByID)
	protectedRouter.HandleFunc("POST /outlet/add", app.SalesOutletHandler.AddSalesOutlet)
	protectedRouter.HandleFunc("POST /outlet/update/{id}", app.SalesOutletHandler.UpdateSalesOutlet)
	protectedRouter.HandleFunc("DELETE /outlet/delete/{id}", app.SalesOutletHandler.DeleteSalesOutlet)

	http.HandleFunc("GET /products/outlet/{id}", app.SalesOutletHandler.GetSalesOutletProducts)
	http.HandleFunc("GET /product/outlet/{outletId}/{productId}", app.SalesOutletHandler.GetSalesOutletProductsByProductID)
	protectedRouter.HandleFunc("POST /stock/add", app.SalesOutletHandler.AddStockItem)
	protectedRouter.HandleFunc("POST /stock/update/{outletId}/{productId}/{amount}/{size}", app.SalesOutletHandler.UpdateStockItem)
	protectedRouter.HandleFunc("DELETE /stock/delete/{outletId}/{productId}", app.SalesOutletHandler.DeleteStockItem)

	protectedRouter.HandleFunc("GET /orders", app.OrderHandler.ListAllOrders)
	protectedRouter.HandleFunc("GET /orders/{orderID}", app.OrderHandler.GetOrderInfo)
	protectedRouter.HandleFunc("GET /users/{userID}/orders", app.OrderHandler.ListAllUserOrders)
	protectedRouter.HandleFunc("GET /sales-outlets/{salesOutletID}/orders", app.OrderHandler.ListOrdersBySalesOutlet)
	protectedRouter.HandleFunc("POST /orders", app.OrderHandler.CreateOrder)
	protectedRouter.HandleFunc("PATCH /orders/{orderID}/status", app.OrderHandler.UpdateOrderStatus)
	protectedRouter.HandleFunc("DELETE /orders/{orderID}", app.OrderHandler.DeleteOrder)
	protectedRouter.HandleFunc("DELETE /order-items/{orderItemID}", app.OrderHandler.DeleteOrderItem)

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
