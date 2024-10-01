package server

import (
	"database/sql"
	"food_delivery/config"
	"food_delivery/repository"
	"food_delivery/server/apiIntegration"
	"food_delivery/server/handler"
	"food_delivery/server/middlware"
	"food_delivery/service"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import the pq driver anonymously
)

func StartServer(cfg *config.Config) {

	// dsn := "postgres://" + cfg.DbUsername + ":" + cfg.DbPassword + "@database/" + cfg.DbName + "?sslmode=disable"
	// db, err := sql.Open("postgres", dsn)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	dsn := "postgres://" + cfg.DbUsername + ":" + cfg.DbPassword + "@localhost/" + cfg.DbName + "?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	tokenService := service.NewTokenService(cfg, db)

	cartRepo := repository.NewCartRepository(db)
	cartHandler := handler.NewCartController(tokenService, cartRepo)

	supplierRepo := repository.NewSupplierRepository(db)
	supplierHandler := handler.SupplierHandler{Repo: &supplierRepo}

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewAuthHandler(tokenService, userRepo, cartRepo)

	menuRepo := repository.NewMenuRepository(db)
	menuHandler := handler.MenuHandler{Repo: &menuRepo}

	categoryRepo := repository.NewCategoryRepository(db)
	categoryHandler := handler.NewcategoryController(categoryRepo)

	addressRepo := repository.NewAddressRepository(db)
	addressHandler := handler.NewAdressController(addressRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderHandler := handler.NewOrderController(orderRepo)

	mux := http.NewServeMux()

	// err = apiIntegration.InsertSuppliers(db)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	go apiIntegration.UpdatingPrice(db)

	// mux.Handle("GET /suppliers", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(supplierHandler.GetAll), tokenService))

	// mux.HandleFunc("GET /suppliers", supplierHandler.GetAll)
	mux.Handle("GET /suppliers", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(supplierHandler.GetAll), tokenService))

	mux.HandleFunc("POST /supplier", supplierHandler.Create)
	mux.HandleFunc("GET /supplier/{id}", supplierHandler.GetbyId)
	mux.HandleFunc("GET /supplier/{id}/menu", menuHandler.GetAll)
	mux.HandleFunc("GET /menu/category/{id}", menuHandler.GetMenubyCategory)

	// mux.HandleFunc("GET /categories", categoryHandler.GetAll)
	mux.Handle("GET /categories", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(categoryHandler.GetAll), tokenService))

	mux.Handle("POST /refresh", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(userHandler.ValidRefreshToken), tokenService))

	// mux.HandleFunc("POST /refresh", userHandler.ValidRefreshToken)

	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.Handle("GET /logout", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(userHandler.Logout), tokenService))

	mux.Handle("GET /cart/createCart", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.Create), tokenService))
	mux.Handle("POST /cart/additem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.AddItemToCart), tokenService))
	mux.Handle("POST /cart/updateCartItem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.UpdateCartItem), tokenService))
	mux.Handle("POST /cart/removeItem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.RemoveItemFromCart), tokenService))
	mux.HandleFunc("GET /cart/getCart", cartHandler.GetCart)
	// mux.Handle("GET /cart/getCart", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.GetCart), tokenService))
	mux.Handle("POST /cart/checkout", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.CheckoutCart), tokenService))

	mux.Handle("POST /address", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(addressHandler.Create), tokenService))
	mux.Handle("GET /address", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(addressHandler.GetAddress), tokenService))

	mux.Handle("GET /orders", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(orderHandler.GetOrders), tokenService))
	mux.Handle("GET /orders/{orderId}", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(orderHandler.GetOrderDetails), tokenService))

	srv := &http.Server{
		Handler: middlware.CORSMiddleware(mux),
		Addr:    cfg.Port,
	}

	log.Println("Server is listening on port 8080")

	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem")
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServeTLS error: %v", err)
	}

}
