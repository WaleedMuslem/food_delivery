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

func FetchMenu(supplier_id int) ([]model.Menu, error) {
	// url := "https://foodapi.golang.nixdev.co/suppliers?limit=10&page=1"
	url := fmt.Sprintf("http://foodapi.golang.nixdev.co/suppliers/%d/menu", supplier_id)

	// Create a custom HTTP client that skips SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 20 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching menus: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var menu []model.Menu
	var response respond.MenuRespond

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	menu = response.Menu

	return menu, nil
}

func InsertMenuItem(db *sql.DB, item model.Menu, supplier_id int) {
	// Insert into menu_items
	var menuItemID int

	query1 := `
		WITH ins_supplier_type AS (
		INSERT INTO category (category_name)
			VALUES ($1)
		ON CONFLICT (category_name) DO NOTHING -- If type already exists, do nothing
			RETURNING category_id
		
		),


		get_supplier_type_id AS (
			SELECT category_id FROM ins_supplier_type
			UNION ALL
			SELECT category_id FROM category WHERE category_name = $1
		)
		INSERT INTO products (id, name, price, supplier_id, image, category_id)
		VALUES ($2, $3, $4, $5, $6, (SELECT category_id FROM get_supplier_type_id LIMIT 1))
		RETURNING id;
			`

	err := db.QueryRow(query1, item.Type, item.ID, item.Name, item.Price, supplier_id, item.Image).Scan(&menuItemID)
	if err != nil {
		log.Fatalf("Error inserting menu item: %v", err)
	}
	// query := `INSERT INTO products (id, name, price, supplier_id, image, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	// err := db.QueryRow(query, item.ID, item.Name, item.Price, supplier_id, item.Image, item.Type).Scan(&menuItemID)
	// if err != nil {
	// 	log.Fatalf("Error inserting menu item: %v", err)
	// }

	// Insert ingredients
	for _, ingredient := range item.Ingredients {
		_, err := db.Exec(`INSERT INTO ingredients (product_id, ingredient) VALUES ($1, $2)`, menuItemID, ingredient)
		if err != nil {
			log.Fatalf("Error inserting ingredient: %v", err)
		}
	}

}

func UpdateMenuItemPrice(db *sql.DB, item model.Menu) {
	// Update only the price for the menu item where id matches
	query := `
        UPDATE products 
        SET price = $2
        WHERE id = $1`

	_, err := db.Exec(query, item.ID, item.Price)
	if err != nil {
		log.Fatalf("Error updating price for menu item: %v", err)
	}

	// fmt.Printf("Updated %s price to %.2f\n", item.Name, item.Price)
}
