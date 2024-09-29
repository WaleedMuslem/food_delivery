package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/model"
	"food_delivery/repository"
	"food_delivery/server/middlware"
	"food_delivery/service"
	"io"
	"log"
	"net/http"
)

type CartHandler struct {
	Repo         repository.ICart
	tokenService *service.TokenService
}

func NewCartController(tokenService *service.TokenService, repo repository.ICart) *CartHandler {
	return &CartHandler{Repo: repo,
		tokenService: tokenService,
	}
}

func (ch *CartHandler) Create(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cartId, err := ch.Repo.CreateCart(claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(cartId)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

func (ch *CartHandler) AddItemToCart(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	// if !ok {
	// 	// Handle the case where the type assertion fails
	// 	log.Print("Failed to retrieve JWT claims from context")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var cartItem model.CartItem

	err = json.Unmarshal(body, &cartItem)
	if err != nil {
		panic(err)
	}

	err = ch.Repo.AddItemToCart(cartItem.CartId, cartItem.ProductID, cartItem.Quantity, cartItem.Price)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func (ch *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	// if !ok {
	// 	// Handle the case where the type assertion fails
	// 	log.Print("Failed to retrieve JWT claims from context")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var cartItem model.CartItem

	err = json.Unmarshal(body, &cartItem)
	if err != nil {
		panic(err)
	}

	err = ch.Repo.UpdateCartItem(cartItem.CartId, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func (ch *CartHandler) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	// if !ok {
	// 	// Handle the case where the type assertion fails
	// 	log.Print("Failed to retrieve JWT claims from context")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var cartItem model.CartItem

	err = json.Unmarshal(body, &cartItem)
	if err != nil {
		panic(err)
	}

	err = ch.Repo.RemoveItemFromCart(cartItem.CartId, cartItem.ProductID)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)

}

func (ch *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {

	// // // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	// if !ok {
	// 	// Handle the case where the type assertion fails
	// 	log.Print("Failed to retrieve JWT claims from context")
	// 	http.Error(w, "Unauthorized to test", http.StatusUnauthorized)
	// 	return
	// }

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 	return
	// }
	// defer r.Body.Close()

	// var cartItem model.CartItem

	// err = json.Unmarshal(body, &cartItem)
	// if err != nil {
	// 	panic(err)
	// }

	claims, err := ch.tokenService.ValidateAccessToken(ch.tokenService.GetTokenFromBearerString(r.Header.Get("Authorization")))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	cart, err := ch.Repo.GetCart(claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(cart)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

func (ch *CartHandler) CheckoutCart(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var cartItem model.CartItem

	err = json.Unmarshal(body, &cartItem)
	if err != nil {
		panic(err)
	}

	cartId, err := ch.Repo.CheckoutCart(cartItem.CartId, claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(cartId)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

// func (ch *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
// 	// Get user ID and product ID from request (e.g., path params)
// 	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
// 	// if !ok {
// 	// 	// Handle the case where the type assertion fails
// 	// 	log.Print("Failed to retrieve JWT claims from context")
// 	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 	// 	return
// 	// }

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var cartItem model.CartItem

// 	err = json.Unmarshal(body, &cartItem)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Remove the product from the cart in the database
// 	err = ch.Repo.RemoveProductFromCart(cartItem.ProductID, cartItem.CartId)
// 	if err != nil {
// 		http.Error(w, "Error removing product from cart", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
