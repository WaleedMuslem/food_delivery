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
