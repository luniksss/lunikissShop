package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"lunikissShop/pkg/app"
)

func main() {
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/lunishop")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS lunishop")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database created successfully!")

	newApp := app.NewApp(db)
	setupRoutes(newApp)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes(app *app.App) {
	http.HandleFunc("GET /product/list", app.ProductHandler.GetProducts)
	http.HandleFunc("GET /product/{id}", app.ProductHandler.GetProductByID)
	http.HandleFunc("POST /product/add", app.ProductHandler.AddProduct)
	http.HandleFunc("POST /product/update/{id}", app.ProductHandler.UpdateProduct)
	http.HandleFunc("DELETE /product/delete/{id}", app.ProductHandler.DeleteProduct)

	http.HandleFunc("GET /outlet/list", app.SalesOutletHandler.GetSalesOutlet)
	http.HandleFunc("GET /outlet/{id}", app.SalesOutletHandler.GetSalesOutletByID)
	http.HandleFunc("POST /outlet/add", app.SalesOutletHandler.AddSalesOutlet)
	http.HandleFunc("POST /outlet/update/{id}", app.SalesOutletHandler.UpdateSalesOutlet)
	http.HandleFunc("DELETE /outlet/delete/{id}", app.SalesOutletHandler.DeleteSalesOutlet)

	http.HandleFunc("GET /products/outlet/{id}", app.SalesOutletHandler.GetSalesOutletProducts)
	http.HandleFunc("GET /product/outlet/{outletId}/{productId}", app.SalesOutletHandler.GetSalesOutletProductsByProductID)
	http.HandleFunc("POST /stock/add", app.SalesOutletHandler.AddStockItem)
	http.HandleFunc("POST /stock/update/{outletId}/{productId}/{amount}/{size}", app.SalesOutletHandler.UpdateStockItem)
	http.HandleFunc("DELETE /stock/delete/{outletId}/{productId}", app.SalesOutletHandler.DeleteStockItem)

	http.HandleFunc("GET /orders", app.OrderHandler.ListAllOrders)
	http.HandleFunc("GET /orders/{orderID}", app.OrderHandler.GetOrderInfo)
	http.HandleFunc("GET /users/{userID}/orders", app.OrderHandler.ListAllUserOrders)
	http.HandleFunc("GET /sales-outlets/{salesOutletID}/orders", app.OrderHandler.ListOrdersBySalesOutlet)
	http.HandleFunc("POST /orders", app.OrderHandler.CreateOrder)
	http.HandleFunc("PATCH /orders/{orderID}/status", app.OrderHandler.UpdateOrderStatus)
	http.HandleFunc("DELETE /orders/{orderID}", app.OrderHandler.DeleteOrder)
	http.HandleFunc("DELETE /order-items/{orderItemID}", app.OrderHandler.DeleteOrderItem)

	http.HandleFunc("GET /api/v1/users", app.UserHandler.ListAllUsers)
	http.HandleFunc("GET /api/v1/users/{id}", app.UserHandler.GetUserByID)
	http.HandleFunc("POST /api/v1/users/by-email", app.UserHandler.GetUserByEmail)

	http.HandleFunc("POST /api/v1/users", app.UserHandler.AddUser)
	http.HandleFunc("PUT /api/v1/users", app.UserHandler.UpdateUser)
	http.HandleFunc("PATCH /api/v1/users/{id}/role", app.UserHandler.UpdateUserRole)
	http.HandleFunc("DELETE /api/v1/users/{id}", app.UserHandler.DeleteUser)
}
