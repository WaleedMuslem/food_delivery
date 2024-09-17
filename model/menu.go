package model

type Menu struct {
	ID          int
	ExtID       int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	SupplierID  int      `json:"supplier_id"`
	Image       string   `json:"image"`
	Type        string   `json:"type"`
	Ingredients []string `json:"ingredients"`
}

type Product struct {
	ID          int
	Name        string
	Price       float64
	SupplierID  int
	CategoryID  int
	Image       string
	Ingredients []string
}
