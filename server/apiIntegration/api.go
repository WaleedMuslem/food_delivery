package apiIntegration

import (
	"database/sql"
	"fmt"
	"food_delivery/model"
	"food_delivery/server/pool"
	"food_delivery/service"
	"log"
	"sync"
	"time"
)

func UpdatingPrice(db *sql.DB) {
	result := make([]model.Menu, 0)
	resErr := make([]error, 0)

	resCh := make(chan any)
	errCh := make(chan error)

	wPool := pool.NewWorkerPool(resCh, errCh).WithBrokerCount(3)

	const updatedTimeByMins = 1
	ticker := time.NewTicker(updatedTimeByMins * time.Minute)
	defer ticker.Stop()

	// Ensure proper synchronization for shared data
	var mu sync.Mutex

	// Goroutine for result collection
	go func() {
		for res := range resCh {
			menuItems := res.([]model.Menu)
			mu.Lock()
			result = append(result, menuItems...)
			mu.Unlock()
		}
	}()

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

		// Append tasks for each supplier
		for _, supplier := range suppliers {
			supplierID := supplier.ID // avoid closure issues
			wPool.Append(func() ([]model.Menu, error) {
				return service.FetchMenu(supplierID)
			})
		}

		// Shutdown the worker pool after appending all tasks
		wPool.Shutdown()

		// Wait for all tasks to complete before proceeding
		time.Sleep(1 * time.Second) // optional wait for pool cleanup

		// Processing fetched menu items
		for _, item := range result {
			// Check if menu item exists
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
			if err != nil {
				log.Fatalf("Error checking if menu item exists: %v", err)
			}

			// Insert or update the menu item
			if !exists {
				service.InsertMenuItem(db, item, item.SupplierID)
			} else {
				service.UpdateMenuItemPrice(db, item)
			}
		}

		// Reset the result and error slices
		result = make([]model.Menu, 0)
		resErr = make([]error, 0)

		fmt.Println("Updated")
	}
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
