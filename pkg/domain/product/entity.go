package product

type SalesOutlet struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type StockItem struct {
	SalesOutletID string `json:"sales_outlet_id"`
	Product       Product
	Size          int `json:"size"`
	Amount        int `json:"amount"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Images      []Image
}

type Image struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	ImagePath string `json:"image_path"`
	OrderNum  int    `json:"order_num"`
}
