package respond

import "food_delivery/model"

type LoginRespond struct {
	AccessToken string `json:"access_token"`
	// RefreshToken string `json:"refresh_token"`
}

type UserRespond struct {
	ID        uint
	Email     string
	FirstName string
	LastName  string
}

type MenuRespond struct {
	Menu []model.Menu
}

type ItemRespond struct {
	ID          int
	ExtID       int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	SupplierID  int      `json:"supplier_id"`
	Image       string   `json:"image"`
	Type        string   `json:"type"`
	SuppierName string   `json:"supplier_name"`
	Ingredients []string `json:"ingredients"`
}
