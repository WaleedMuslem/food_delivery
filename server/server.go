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

	dsn := "postgres://" + cfg.DbUsername + ":" + cfg.DbPassword + "@localhost/" + cfg.DbName + "?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	supplierRepo := repository.NewSupplierRepository(db)
	supplierHandler := handler.SupplierHandler{Repo: &supplierRepo}

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService(cfg, db)
	userHandler := handler.NewAuthHandler(tokenService, userRepo)

	menuRepo := repository.NewMenuRepository(db)
	menuHandler := handler.MenuHandler{Repo: &menuRepo}

	categoryRepo := repository.NewCategoryRepository(db)
	categoryHandler := handler.NewcategoryController(categoryRepo)

	cartRepo := repository.NewCartRepository(db)
	cartHandler := handler.NewCartController(cartRepo)

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

	mux.HandleFunc("POST /refresh", userHandler.ValidRefreshToken)

	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.Handle("GET /logout", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(userHandler.Logout), tokenService))

	mux.Handle("GET /cart/createCart", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.Create), tokenService))
	mux.Handle("GET /cart/additem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.AddItemToCart), tokenService))
	mux.Handle("GET /cart/updateCartItem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.UpdateCartItem), tokenService))
	mux.Handle("POST /cart/removeItem", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.RemoveItemFromCart), tokenService))
	mux.Handle("GET /cart/getCart", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.GetCart), tokenService))
	mux.Handle("GET /cart/checkout", middlware.AcessTokenValdityMiddleware(http.HandlerFunc(cartHandler.CheckoutCart), tokenService))

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
