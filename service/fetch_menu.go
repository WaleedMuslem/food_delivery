package service

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"food_delivery/model"
	"food_delivery/respond"
	"io"
	"log"
	"net/http"
	"time"
)

func FetchAndUpdateMenu(db *sql.DB, supplier_id int) error {
	// url := "https://foodapi.golang.nixdev.co/suppliers?limit=10&page=1"
	url := fmt.Sprintf("http://foodapi.golang.nixdev.co/suppliers/%d/menu", supplier_id)

	// Create a custom HTTP client that skips SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 20 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching menus: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error . response body: %v", err)
	}

	var menu []model.Menu
	var response respond.MenuRespond

	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	menu = response.Menu

	// Assign the supplier ID to each menu item
	for _, item := range menu {
		item.SupplierID = supplier_id
		// fmt.Println(item.ExtID)
		// var exists bool
		// err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
		// if err != nil {
		// 	log.Fatalf("Error checking if menu item exists: %v", err)
		// }

		// // Insert or update the menu item
		// if !exists {
		UpdateMenuItem(db, item)
		// } else {
		// UpdateMenuItemPrice(db, item)
		// }
	}

	return nil
}

func UpdateMenuItem(db *sql.DB, item model.Menu) {
	// Insert into menu_items
	var menuItemID int

	// fmt.Println(item.SupplierID)
	query := `WITH ins_category_type AS (
		INSERT INTO category (category_name)
		VALUES ($1)
		ON CONFLICT (category_name) DO NOTHING -- If category already exists, do nothing
		RETURNING category_id
	),

	get_category_type_id AS (
		SELECT category_id FROM ins_category_type
		UNION ALL
		SELECT category_id FROM category WHERE category_name = $1
	)

	INSERT INTO products (name, price, supplier_id, image, category_id, ext_id)
	VALUES ($2, $3, $4, $5, (SELECT category_id FROM get_category_type_id LIMIT 1), $6)
	ON CONFLICT (ext_id) DO UPDATE
	SET name = excluded.name,
		price = excluded.price,
		supplier_id = excluded.supplier_id,
		image = excluded.image,
		category_id = (SELECT category_id FROM get_category_type_id LIMIT 1),
		ext_id = excluded.ext_id
	RETURNING id;
	`
	IngredientsItemQuery := `
        INSERT INTO product_ingredient (product_id, ingredient_id)
        VALUES ($1, $2)
        ON CONFLICT (product_id, ingredient_id) DO NOTHING;`

	// 	INSERT INTO suppliers (name, image, opening, closing, ext_id, type_id)
	// VALUES ($2, $3, $4, $5, $6, (SELECT id FROM type_id))
	// ON CONFLICT (ext_id) DO NOTHING;`,

	err := db.QueryRow(query, item.Type, item.Name, item.Price, item.SupplierID, item.Image, item.ExtID).Scan(&menuItemID)
	if err != nil {
		log.Fatalf("Error inserting menu item: %v", err)
	}

	for _, ingredient := range item.Ingredients {
		var ingredientID int
		err := db.QueryRow(`INSERT INTO ingredients (ingredient) VALUES ($1) ON CONFLICT (ingredient) DO UPDATE SET ingredient = EXCLUDED.ingredient RETURNING id;`, ingredient).Scan(&ingredientID)
		if err != nil {
			log.Fatalf("Error inserting ingredient: %v", err)
		}

		_, err = db.Exec(IngredientsItemQuery, menuItemID, ingredientID)
		if err != nil {
			log.Fatalf("Error inserting ingredient_item: %v", err)
		}
	}

	// for _, ingredient := range item.Ingredients {
	// 	_, err := db.Exec(`INSERT INTO ingredients (product_id, ingredient) VALUES ($1, $2)`, menuItemID, ingredient)
	// 	if err != nil {
	// 		log.Fatalf("Error inserting ingredient: %v", err)
	// 	}
	// }

}

// func UpdateMenuItemPrice(db *sql.DB, item model.Menu) {
// 	// Update only the price for the menu item where id matches
// 	query := `
//         UPDATE products
//         SET price = $2
//         WHERE id = $1`

// 	_, err := db.Exec(query, item.ID, item.Price)
// 	if err != nil {
// 		log.Fatalf("Error updating price for menu item: %v", err)
// 	}

// 	// fmt.Printf("Updated %s price to %.2f\n", item.Name, item.Price)
// }
