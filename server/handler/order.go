package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/repository"
	"food_delivery/server/middlware"
	"food_delivery/service"
	"log"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	Repo *repository.OrderRepository
}

func NewOrderController(repo repository.OrderRepository) OrderHandler {
	return OrderHandler{Repo: &repo}
}

func (oh *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := oh.Repo.GetUserOrders(claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

// Handler function to get order details
func (oh *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idString := r.PathValue("orderId")
	// fmt.Println(idString)
	orderId, _ := strconv.Atoi(idString)

	orderDetail, err := oh.Repo.GetOrderDetails(orderId, claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(orderDetail)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}
