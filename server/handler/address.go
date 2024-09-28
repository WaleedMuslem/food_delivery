package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/model"
	"food_delivery/repository"
	"food_delivery/server/middlware"
	"food_delivery/service"
	"log"
	"net/http"
)

type AddressHandler struct {
	Repo *repository.AddressRepository
}

func NewAdressController(repo repository.AddressRepository) AddressHandler {
	return AddressHandler{Repo: &repo}
}

func (ah AddressHandler) Create(w http.ResponseWriter, r *http.Request) {

	// // Type assertion to convert from `any` to `*service.JwtCustomClaims`
	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req := new(model.Address)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Save the address in the database
	addressId, err := ah.Repo.Create(req, claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(addressId)

}

func (ah AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	if !ok {
		// Handle the case where the type assertion fails
		log.Print("Failed to retrieve JWT claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	address, err := ah.Repo.Get(claims.ID)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(address)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}
