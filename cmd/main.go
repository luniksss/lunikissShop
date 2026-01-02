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
}
