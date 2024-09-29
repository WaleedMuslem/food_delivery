package repository

import (
	"errors"
	"food_delivery/model"
)

type CartRepositoryFake struct {
	carts []*model.Cart
}

func NewCartRepositoryFake() ICart {

	carts := []*model.Cart{
		{
			CartID: 1,
			Items: []model.CartItem{
				{
					CartId:     1,
					Name:       "pizza",
					Image:      "image",
					ProductID:  1,
					Quantity:   1,
					Price:      20,
					TotalPrice: 20,
				},
				// {
				// 	CartId:     2,
				// 	Name:       "pizza2",
				// 	Image:      "image2",
				// 	ProductID:  2,
				// 	Quantity:   2,
				// 	Price:      30,
				// 	TotalPrice: 60,
				// },
			},
		},
	}

	return &CartRepositoryFake{
		carts: carts,
	}
}

func (fake CartRepositoryFake) GetCart(userID uint) (*model.Cart, error) {
	if userID == uint(fake.carts[0].CartID) {
		return fake.carts[0], nil
	}

	return nil, errors.New("user not found")
}

func (fake CartRepositoryFake) CheckoutCart(cartID int, userID uint) (int, error) {
	return 0, nil
}

func (fake CartRepositoryFake) CreateCart(userID uint) (int, error) {
	return 0, nil
}
func (fake CartRepositoryFake) AddItemToCart(cartID int, productID int, quantity int, price float64) error {
	return nil
}
func (fake CartRepositoryFake) UpdateCartItem(cartID int, productID int, newQuantity int) error {
	return nil
}
func (fake CartRepositoryFake) RemoveItemFromCart(cartID int, productID int) error {
	return nil
}
