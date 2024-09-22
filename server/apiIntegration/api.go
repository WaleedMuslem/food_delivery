package apiIntegration

import (
	"database/sql"
	"fmt"
	"food_delivery/model"
	"food_delivery/server/pool"
	"food_delivery/service"
	"net/http"
	"sync"
	"time"
)

func UpdatingPrice(db *sql.DB) {
	resErr := make([]error, 0)

	errCh := make(chan error)

	wPool := pool.NewWorkerPool(errCh).WithBrokerCount(3)

	const updatedTimeByMins = 1
	ticker := time.NewTicker(updatedTimeByMins * time.Minute)
	defer ticker.Stop()

	// Ensure proper synchronization for shared data
	var mu sync.Mutex

	// Goroutine for error collection
	go func() {
		for err := range errCh {
			mu.Lock()
			resErr = append(resErr, err)
			mu.Unlock()
		}
	}()

	for range ticker.C {
		// Start the worker pool at the start of each ticker
		wPool.Start()

		// Fetch all suppliers
		suppliers, err := service.FetchAllSuppliers()
		if err != nil {
			fmt.Println(err)
			return
		}
		// store the suppliers into our database
		err = CreateSuppliers(suppliers, db)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Append tasks for each supplier
		for _, supplier := range suppliers {
			supplierID := supplier.ExtID
			wPool.Append(func() error {
				return service.FetchAndUpdateMenu(db, supplierID)
			})
		}

		// Shutdown the worker pool after appending all tasks
		wPool.Shutdown()

		// Wait for all tasks to complete before proceeding
		time.Sleep(1 * time.Second)

		// Processing fetched menu items
		// for _, item := range result {
		// 	// Check if menu item exists
		// 	var exists bool
		// 	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
		// 	if err != nil {
		// 		log.Fatalf("Error checking if menu item exists: %v", err)
		// 	}

		// 	// Insert or update the menu item
		// 	if !exists {
		// 		service.InsertMenuItem(db, item, item.SupplierID)
		// 	} else {
		// 		service.UpdateMenuItemPrice(db, item)
		// 	}
		// }

		// Reset the result and error slices
		resErr = make([]error, 0)

		fmt.Println("Updated")
	}
}

// func InsertSuppliers(db *sql.DB) error {

// }

var defaultImagesByType = map[string]string{
	"restaurant":  "https://images.pexels.com/photos/67468/pexels-photo-67468.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
	"bar":         "https://images.pexels.com/photos/941864/pexels-photo-941864.jpeg",
	"supermarket": "https://images.pexels.com/photos/264636/pexels-photo-264636.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
	"coffee_shop": "https://images.pexels.com/photos/28543481/pexels-photo-28543481/free-photo-of-cozy-london-cafe-with-coffee-in-blue-cup.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
	"shop":        "https://images.pexels.com/photos/135620/pexels-photo-135620.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
	"Other":       "https://images.pexels.com/photos/1640777/pexels-photo-1640777.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1", // Default for any other type
}

func CreateSuppliers(suppliers []model.SupplierFromAPI, db *sql.DB) error {
	var finalErr error

	// 	query := `
	// WITH ins AS (
	//     INSERT INTO supplier_types (type)
	//     VALUES ($1)
	//     ON CONFLICT (type) DO NOTHING
	//     RETURNING id, type
	// )
	// INSERT INTO suppliers (name, image, opening, closing, ext_id ,type_id)
	// VALUES ($2, $3, $4, $5, $6, COALESCE((SELECT id FROM supplier_types WHERE type = $1), (SELECT id FROM ins WHERE type = $1)) ON CONFLICT (ext_id) DO NOTHING;`

	// `
	// WITH ins AS (
	// 	INSERT INTO supplier_types (type)
	// 	VALUES ($1)
	// 	ON CONFLICT (type) DO NOTHING
	// 	RETURNING id, type
	// )
	// INSERT INTO suppliers (name, image, opening, closing, ext_id ,type_id)
	// VALUES ($2, $3, $4, $5, $6, COALESCE((SELECT id FROM supplier_types WHERE type = $1), (SELECT id FROM ins WHERE type = $1)) ON CONFLICT (ext_id) DO NOTHING;`

	// Default images for different supplier types

	for _, supplier := range suppliers {

		if !isValidImage2(supplier.Image) {
			// Use a type-based placeholder image
			switch supplier.Type {
			case "restaurant":
				supplier.Image = defaultImagesByType["restaurant"]
			case "bar":
				supplier.Image = defaultImagesByType["bar"]
			case "supermarket":
				supplier.Image = defaultImagesByType["supermarket"]
			case "coffee_shop":
				supplier.Image = defaultImagesByType["coffee_shop"]
			case "shop":
				supplier.Image = defaultImagesByType["shop"]
			default:
				supplier.Image = defaultImagesByType["Other"]

			}

			// if placeholder, exists := defaultImagesByType[supplier.Type]; exists {

			// 	supplier.Image = placeholder
			// } else {
			// 	// Use "Other" placeholder if type doesn't match any known types
			// 	supplier.Image = defaultImagesByType["Other"]
			// }
		}

		_, err := db.Exec(`
		WITH ins AS (
			INSERT INTO supplier_type (type)
			VALUES ($1)
			ON CONFLICT (type) DO NOTHING
			RETURNING id
		),
		type_id AS (
			SELECT id FROM supplier_type WHERE type = $1
			UNION
			SELECT id FROM ins
		) 
		INSERT INTO suppliers (name, image, opening, closing, ext_id, type_id)
		VALUES ($2, $3, $4, $5, $6, (SELECT id FROM type_id))
		ON CONFLICT (ext_id) DO UPDATE
		SET name = EXCLUDED.name,
			image = EXCLUDED.image,
			opening = EXCLUDED.opening,
			closing = EXCLUDED.closing,
			type_id = (SELECT id FROM type_id);`,
			supplier.Type, supplier.Name, supplier.Image, supplier.WorkingHours.Opening, supplier.WorkingHours.Closing, supplier.ExtID,
		)

		if err != nil {
			finalErr = fmt.Errorf("error processing supplier %s: %w", supplier.Name, err)
		}
	}

	return finalErr
}

func isValidImage2(imageURL string) bool {
	if imageURL == "" {
		return false
	}

	// Send a HEAD request to check if the image URL exists and is accessible
	resp, err := http.Head(imageURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	// Optional: Check for image content type (e.g., "image/png", "image/jpeg")
	contentType := resp.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return false
	}

	return true
}

// func UpdatingPrice(db *sql.DB) {

// 	result := make([]model.Menu, 0)
// 	resErr := make([]error, 0)

// 	resCh := make(chan any)
// 	errCh := make(chan error)

// 	wPool := pool.NewWorkerPool(resCh, errCh).WithBrokerCount(3)
// 	wPool.Start()

// 	const updatedTimeByMins = 1
// 	ticker := time.NewTicker(updatedTimeByMins * time.Minute)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		// Call the function every time the ticker ticks
// 		suppliers, err := service.FetchAllSuppliers()
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		for _, supplier := range suppliers {

// 			wPool.Append(func() ([]model.Menu, error) {
// 				return service.FetchMenu(supplier.ID)
// 			})

// 			go func() {
// 				for res := range resCh {
// 					result = append(result, res.(model.Menu))
// 				}
// 			}()

// 			go func() {
// 				for err := range errCh {
// 					resErr = append(resErr, err)
// 				}
// 			}()

// 			wPool.Shutdown()
// 			fmt.Println(result)
// 			fmt.Println(resErr)
// 			// menu, err := service.FetchMenu(supplier.ID)
// 			// if err != nil {
// 			// 	fmt.Println(err)
// 			// 	continue
// 			// }
// 			for _, item := range result {

// 				// Check if menu item exists
// 				var exists bool
// 				err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
// 				if err != nil {
// 					log.Fatalf("Error checking if menu item exists: %v", err)
// 				}

// 				if !exists {
// 					// Insert a new menu item instead of updating
// 					service.InsertMenuItem(db, item, supplier.ID)
// 				} else {
// 					// Proceed with updating the menu item
// 					service.UpdateMenuItemPrice(db, item)
// 				}

// 			}

// 		}

// 		fmt.Println("Updated")
// 	}
// }

// func UpdatingPrice(db *sql.DB) {

// 	result := make([]string, 0)
// 	resErr := make([]error, 0)

// 	resCh := make(chan any)
// 	errCh := make(chan error)

// 	wPool := pool.NewWorkerPool(resCh, errCh).WithBrokerCount(3)
// 	wPool.Start()

// 	const updatedTimeByMins = 1
// 	ticker := time.NewTicker(updatedTimeByMins * time.Minute)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		// Call the function every time the ticker ticks
// 		suppliers, err := service.FetchAllSuppliers()
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		for _, supplier := range suppliers {
// 			menu, err := service.FetchMenu(supplier.ID)
// 			if err != nil {
// 				fmt.Println(err)
// 				continue
// 			}
// 			for _, item := range menu {

// 				// Check if menu item exists
// 				var exists bool
// 				err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
// 				if err != nil {
// 					log.Fatalf("Error checking if menu item exists: %v", err)
// 				}

// 				if !exists {
// 					// Insert a new menu item instead of updating
// 					service.InsertMenuItem(db, item, supplier.ID)
// 				} else {
// 					// Proceed with updating the menu item
// 					service.UpdateMenuItemPrice(db, item)
// 				}

// 			}

// 		}

// 		fmt.Println("Updated")
// 	}
// }
