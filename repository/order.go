package repository

import (
	"database/sql"
	"food_delivery/model"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return OrderRepository{Db: db}
}

func (or OrderRepository) GetUserOrders(userID uint) ([]model.UserOrder, error) {
	rows, err := or.Db.Query(`
        SELECT o.order_id, o.total_amount, o.created_at, o.status,
               p.id, p.name AS product_name, op.quantity, op.total_price
        FROM orders o
        JOIN order_product op ON o.order_id = op.order_id
        JOIN products p ON op.product_id = p.id
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*model.UserOrder) // Map to store orders
	for rows.Next() {
		var product model.OrderProduct
		var order model.UserOrder

		// Scan data from the query
		err := rows.Scan(&order.OrderID, &order.TotalAmount, &order.CreatedAt, &order.Status,
			&product.ProductID, &product.ProductName, &product.Quantity, &product.TotalPrice)
		if err != nil {
			return nil, err
		}

		// Check if the order already exists in the map
		if existingOrder, found := ordersMap[order.OrderID]; found {
			// Append the product to the existing order's products list
			existingOrder.Products = append(existingOrder.Products, product)
		} else {
			// Create a new order entry
			order.Products = []model.OrderProduct{product}
			ordersMap[order.OrderID] = &order
		}
	}

	// Convert the map to a slice of orders
	var orders []model.UserOrder
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (or OrderRepository) GetOrderDetails(orderId int, userID uint) ([]model.UserOrder, error) {

	rows, err := or.Db.Query(`
        SELECT o.order_id, o.total_amount, o.created_at, o.status,
               p.id, p.name AS product_name, op.quantity, op.total_price
        FROM orders o
        JOIN order_product op ON o.order_id = op.order_id
        JOIN products p ON op.product_id = p.id
        WHERE o.user_id = $1 AND o.order_id = $2
        ORDER BY o.created_at DESC`, userID, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*model.UserOrder) // Map to store orders
	for rows.Next() {
		var product model.OrderProduct
		var order model.UserOrder

		// Scan data from the query
		err := rows.Scan(&order.OrderID, &order.TotalAmount, &order.CreatedAt, &order.Status,
			&product.ProductID, &product.ProductName, &product.Quantity, &product.TotalPrice)
		if err != nil {
			return nil, err
		}

		// Check if the order already exists in the map
		if existingOrder, found := ordersMap[order.OrderID]; found {
			// Append the product to the existing order's products list
			existingOrder.Products = append(existingOrder.Products, product)
		} else {
			// Create a new order entry
			order.Products = []model.OrderProduct{product}
			ordersMap[order.OrderID] = &order
		}
	}

	// Convert the map to a slice of orders
	var orders []model.UserOrder
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return orders, nil
	// Query the order details, including address data
	// 	orderQuery := `
	// 		SELECT o.order_id, o.created_at, o.status,
	// 			a.city, a.street, a.zip_code, a.floor, a.apartment,
	// 			o.total_amount,
	// 			p.id AS product_id, p.name AS product_name, op.quantity, op.total_price
	// 		FROM orders o
	// 		JOIN order_product op ON o.order_id = op.order_id
	// 		JOIN products p ON op.product_id = p.id
	// 		JOIN addresses a ON o.user_id = a.user_id -- Assuming the address table links directly to users
	// 		WHERE o.order_id = $1

	//  `

	// rows, err := or.Db.Query(orderQuery, orderId)
	// if err != nil {
	// 	return model.OrderDetails{}, err
	// }
	// defer rows.Close()

	// // Initialize the response structure
	// var orderDetails model.OrderDetails
	// var products []model.OrderProduct

	// // Loop through the rows and scan the data into the order details and product structures
	// for rows.Next() {
	// 	var product model.OrderProduct
	// 	err := rows.Scan(&orderDetails.OrderID, &orderDetails.CreatedAt, &orderDetails.Status,
	// 		&orderDetails.DeliveryAddress.City, &orderDetails.DeliveryAddress.Street,
	// 		&orderDetails.DeliveryAddress.Zip, &orderDetails.DeliveryAddress.Floor,
	// 		&orderDetails.DeliveryAddress.Apartment, &orderDetails.TotalAmount,
	// 		&product.ProductID, &product.ProductName, &product.Quantity, &product.TotalPrice)
	// 	if err != nil {
	// 		return model.OrderDetails{}, err
	// 	}
	// 	products = append(products, product)
	// }

	// // Attach the products to the order details
	// orderDetails.Products = products

	// return orderDetails, nil

}
