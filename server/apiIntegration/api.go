package apiIntegration

import (
	"database/sql"
	"fmt"
	"food_delivery/model"
	"food_delivery/server/pool"
	"food_delivery/service"
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

func CreateSuppliers(suppliers []model.Supplier, db *sql.DB) error {
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

	for _, supplier := range suppliers {
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
	ON CONFLICT (ext_id) DO NOTHING;`,
			supplier.Type, supplier.Name, supplier.Image, supplier.WorkingHours.Opening, supplier.WorkingHours.Closing, supplier.ExtID,
		)

		if err != nil {
			finalErr = fmt.Errorf("error processing supplier %s: %w", supplier.Name, err)
		}
	}

	return finalErr
}

// func fetchingandupdateingWithWorker() {

// }

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
