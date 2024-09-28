package model

type OrderProduct struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type UserOrder struct {
	OrderID     int            `json:"order_id"`
	TotalAmount float64        `json:"total_amount"`
	CreatedAt   string         `json:"created_at"`
	Status      string         `json:"status"`
	Products    []OrderProduct `json:"products"`
}
type OrderDetails struct {
	OrderID         int            `json:"order_id"`
	CreatedAt       string         `json:"created_at"`
	Status          string         `json:"status"`
	DeliveryAddress Address        `json:"delivery_address"`
	TotalAmount     float64        `json:"total_amount"`
	Products        []OrderProduct `json:"products"`
}
