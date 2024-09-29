package repository

import (
	"database/sql"
	"food_delivery/model"
)

type ICart interface {
	CreateCart(userID uint) (int, error)
	AddItemToCart(cartID int, productID int, quantity int, price float64) error
	UpdateCartItem(cartID int, productID int, newQuantity int) error
	RemoveItemFromCart(cartID int, productID int) error
	GetCart(userID uint) (*model.Cart, error)
	CheckoutCart(cartID int, userID uint) (int, error)
}

type CartRepository struct {
	Db *sql.DB
}

func NewCartRepository(db *sql.DB) ICart {
	return &CartRepository{Db: db}
}

func (cr CartRepository) CreateCart(userID uint) (int, error) {
	var cartID int

	// Check if the user already has an unordered cart
	err := cr.Db.QueryRow(`SELECT cart_id FROM carts WHERE user_id = $1 AND is_ordered = false`, userID).Scan(&cartID)
	if err == nil {
		// If the cart exists, return the cart ID
		return cartID, nil
	}

	// If no unordered cart exists, create a new one
	err = cr.Db.QueryRow(`INSERT INTO carts (user_id) 
                        VALUES ($1) RETURNING cart_id`, userID).Scan(&cartID)
	if err != nil {
		return 0, err
	}

	// Return the new cart ID
	return cartID, nil
}

// func (cr CartRepository) CreateCart(userID uint) error {

// 	err := cr.Db.QueryRow(`INSERT INTO carts (user_id)
//                         VALUES ($1) RETURNING cart_id`, userID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }

func (cr CartRepository) AddItemToCart(cartID int, productID int, quantity int, price float64) error {
	_, err := cr.Db.Exec(`INSERT INTO cart_product (cart_id, product_id, quantity, price) 
                       VALUES ($1, $2, $3, $4) 
                       ON CONFLICT (cart_id, product_id) 
                       DO UPDATE SET quantity = $3`,
		cartID, productID, quantity, price)
	return err
}

func (cr CartRepository) UpdateCartItem(cartID int, productID int, newQuantity int) error {
	_, err := cr.Db.Exec(`UPDATE cart_product SET quantity = $3 
                       WHERE cart_id = $1 AND product_id = $2`, cartID, productID, newQuantity)
	return err
}

func (cr CartRepository) RemoveItemFromCart(cartID int, productID int) error {
	_, err := cr.Db.Exec(`DELETE FROM cart_product WHERE cart_id = $1 AND product_id = $2`, cartID, productID)
	return err
}

func (cr CartRepository) GetCart(userID uint) (*model.Cart, error) {
	var cart model.Cart
	err := cr.Db.QueryRow(`SELECT cart_id FROM carts WHERE user_id = $1 AND is_ordered = FALSE`, userID).Scan(&cart.CartID)
	if err != nil {
		return &model.Cart{}, err
	}

	rows, err := cr.Db.Query(`
    SELECT 
        cp.cart_id, 
        p.name, 
        p.image, 
        cp.product_id, 
        cp.quantity, 
        cp.price, 
        cp.total_price 
    FROM 
        cart_product cp
    JOIN 
        products p 
    ON 
        cp.product_id = p.id
    WHERE 
        cp.cart_id = $1`, cart.CartID)
	if err != nil {
		return &model.Cart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.CartItem
		err = rows.Scan(&item.CartId, &item.Name, &item.Image, &item.ProductID, &item.Quantity, &item.Price, &item.TotalPrice)
		if err != nil {
			return &model.Cart{}, err
		}
		cart.Items = append(cart.Items, item)
	}
	return &cart, nil
}

func (cr CartRepository) CheckoutCart(cartID int, userID uint) (int, error) {
	tx, err := cr.Db.Begin() // Start a transaction
	if err != nil {
		return 0, err
	}

	// Insert into orders
	_, err = tx.Exec(`INSERT INTO orders (user_id, cart_id, total_amount) 
                      SELECT $2, cart_id, SUM(total_price) FROM cart_product 
                      WHERE cart_id = $1 GROUP BY cart_id`, cartID, userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get the order_id of the newly created order
	var orderID int
	err = tx.QueryRow(`SELECT order_id FROM orders WHERE cart_id = $1`, cartID).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert cart products into order_product
	_, err = tx.Exec(`INSERT INTO order_product (order_id, product_id, quantity, price) 
                      SELECT $1, product_id, quantity, price FROM cart_product 
                      WHERE cart_id = $2`, orderID, cartID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Mark the cart as ordered
	_, err = tx.Exec(`UPDATE carts SET is_ordered = TRUE WHERE cart_id = $1`, cartID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = tx.Exec(`DELETE FROM cart_product
						WHERE cart_id = $1;
						`, cartID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var newCartId int
	// Insert new cart and return the new cart_id
	err = tx.QueryRow(`INSERT INTO carts (user_id) 
                      VALUES ($1) RETURNING cart_id`, userID).Scan(&newCartId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// cartId, err := cr.CreateCart(userID)
	// if err != nil {
	// 	return 0, err
	// }

	return newCartId, tx.Commit() // Commit the transaction
}

// func (cr CartRepository) RemoveProductFromCart(productId int, cartId int) error {

// 	_, err := cr.Db.Exec(`DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`, cartId, productId)

// 	return err

// }
