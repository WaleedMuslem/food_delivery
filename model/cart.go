package model

type Cart struct {
	CartID int
	Items  []CartItem
}

type CartItem struct {
	CartId     int     `json:"cart_id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}
