package repository

import (
	"database/sql"
	"food_delivery/model"
	"food_delivery/respond"

	"github.com/lib/pq"
)

// type Product struct {
// 	ID          int
// 	Name        string
// 	Price       float64
// 	SupplierID  int
// 	CategoryID  int
// 	Image       string
// 	Ingredients []string
// }

type IMenu interface {
	GetAll() ([]model.Menu, error)
}

type MenuRepository struct {
	Db *sql.DB
}

func NewMenuRepository(db *sql.DB) MenuRepository {
	return MenuRepository{Db: db}
}

func (mr MenuRepository) GetAllBySupplierId(supplier_id int) ([]respond.ItemRespond, error) {

	menus := []respond.ItemRespond{}

	query := `
        SELECT p.id, p.name, p.price, p.supplier_id, s.name, p.image, c.category_name, 
               array_agg(i.ingredient) as ingredients
        FROM products p
        JOIN category c ON p.category_id = c.category_id
        JOIN suppliers s ON p.supplier_id = s.ext_id  -- Join with suppliers table
        LEFT JOIN product_ingredient pi ON p.id = pi.product_id
        LEFT JOIN ingredients i ON pi.ingredient_id = i.id
        WHERE s.ext_id = $1
        GROUP BY p.id, s.name, c.category_name  -- Include supplier_name in GROUP BY
    `

	result, err := mr.Db.Query(query, supplier_id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		menu := respond.ItemRespond{}
		var ingredients pq.StringArray

		err := result.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.SupplierID, &menu.SuppierName, &menu.Image, &menu.Type, &ingredients)
		if err != nil {
			return nil, err
		}

		// Assign the array of ingredients to the menu struct
		menu.Ingredients = ingredients

		menus = append(menus, menu)
	}

	return menus, nil
}

func (mr MenuRepository) GetMenuByCategory(category_id int) ([]respond.ItemRespond, error) {
	menus := []respond.ItemRespond{}

	query := `
        SELECT p.id, p.name, p.price, p.supplier_id, s.name, p.image, c.category_name, 
               array_agg(i.ingredient) as ingredients
        FROM products p
        JOIN category c ON p.category_id = c.category_id
        JOIN suppliers s ON p.supplier_id = s.ext_id  -- Join with suppliers table
        LEFT JOIN product_ingredient pi ON p.id = pi.product_id
        LEFT JOIN ingredients i ON pi.ingredient_id = i.id
        WHERE p.category_id = $1
        GROUP BY p.id, s.name, c.category_name  -- Include supplier_name in GROUP BY
    `

	result, err := mr.Db.Query(query, category_id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		menu := respond.ItemRespond{}
		var ingredients pq.StringArray

		err := result.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.SupplierID, &menu.SuppierName, &menu.Image, &menu.Type, &ingredients)
		if err != nil {
			return nil, err
		}

		// Assign the array of ingredients to the menu struct
		menu.Ingredients = ingredients

		menus = append(menus, menu)
	}

	return menus, nil
}

// func (mr MenuRepository) GetMenuByCategory(category_id int) ([]model.Menu, error) {
// 	menus := []model.Menu{}

// 	query := `
//         SELECT p.id, p.name, p.price, p.supplier_id, p.image, c.category_name,
//                array_agg(i.ingredient) as ingredients
//         FROM products p
//         JOIN category c ON p.category_id = c.category_id
//         LEFT JOIN product_ingredient pi ON p.id = pi.product_id
//         LEFT JOIN ingredients i ON pi.ingredient_id = i.id
//         WHERE p.category_id = $1
//         GROUP BY p.id, c.category_name
//     `

// 	result, err := mr.Db.Query(query, category_id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer result.Close()

// 	for result.Next() {
// 		menu := model.Menu{}
// 		var ingredients pq.StringArray

// 		err := result.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.SupplierID, &menu.Image, &menu.Type, &ingredients)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Assign the array of ingredients to the menu struct
// 		menu.Ingredients = ingredients

// 		menus = append(menus, menu)
// 	}

// 	return menus, nil
// }

// func (sr SupplierRepository) GetSyppliersBycategoryId(id int) ([]model.Supplier, error) {

// 	suppliersByCategory := []model.Supplier{}

// 	// Updated query to join with the type table
// 	result, err := sr.Db.Query(`
// 		SELECT s.id, s.name, s.image, s.opening, s.closing, s.ext_id, s.type_id, t.type
// 		FROM suppliers s
// 		JOIN supplier_type t ON s.type_id = t.id
// 		Where categ
// 		ORDER BY s.id DESC
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer result.Close()

// 	for result.Next() {
// 		supplier := model.Supplier{}
// 		supplierType := model.Type{} // Assuming there's a Type struct in model

// 		// Scanning both supplier and type data
// 		err := result.Scan(
// 			&supplier.ID,
// 			&supplier.Name,
// 			&supplier.Image,
// 			&supplier.WorkingHours.Opening,
// 			&supplier.WorkingHours.Closing,
// 			&supplier.ExtID,
// 			&supplierType.ID,   // Scanning type_id into Type struct's ID
// 			&supplierType.Type, // Scanning the type field
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Assign the scanned type to the supplier's Type field
// 		supplier.Type = supplierType

// 		suppliersByCategory = append(suppliersByCategory, supplier)
// 	}

// 	return suppliersByCategory, nil
// }
