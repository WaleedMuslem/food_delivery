package apiIntegration

import (
	"database/sql"
	"fmt"
	"food_delivery/service"
	"log"
	"time"
)

func UpdatingPrice(db *sql.DB) {

	const updatedTimeByMins = 1
	ticker := time.NewTicker(updatedTimeByMins * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Call the function every time the ticker ticks
		suppliers, err := service.FetchAllSuppliers()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, supplier := range suppliers {
			menu, err := service.FetchMenu(supplier.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, item := range menu {

				// Check if menu item exists
				var exists bool
				err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", item.ID).Scan(&exists)
				if err != nil {
					log.Fatalf("Error checking if menu item exists: %v", err)
				}

				if !exists {
					// Insert a new menu item instead of updating
					service.InsertMenuItem(db, item, supplier.ID)
				} else {
					// Proceed with updating the menu item
					service.UpdateMenuItemPrice(db, item)
				}

			}

		}

		fmt.Println("Updated")
	}
}
