package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/model"
	"food_delivery/repository"
	"io"
	"net/http"
	"strconv"
)

type SupplierHandler struct {
	Repo *repository.SupplierRepository
}

func NewSupplierController(repo repository.SupplierRepository) SupplierHandler {
	return SupplierHandler{Repo: &repo}
}

func (sh *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var supplier model.Supplier

	err = json.Unmarshal(body, &supplier)
	if err != nil {
		panic(err)
	}

	err = sh.Repo.Create(supplier)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(supplier)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

func (sh *SupplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	// Type assertion to convert from `any` to `*service.JwtCustomClaims`
	// claims, ok := r.Context().Value(middlware.ClaimsKey).(*service.JwtCustomClaims)
	// if !ok {
	// 	// Handle the case where the type assertion fails
	// 	log.Print("Failed to retrieve JWT claims from context")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// fmt.Println(claims)

	suppliers, err := sh.Repo.GetAll()
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(suppliers)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	// p1, _ := bcrypt.GenerateFromPassword([]byte("123456789"), bcrypt.MinCost)
	// fmt.Println(string(p1))

	w.WriteHeader(200)
	w.Write(jsonData)
}

func (sh *SupplierHandler) GetbyId(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	// fmt.Println(idString)
	id, _ := strconv.Atoi(idString)

	// fmt.Println(id)

	supplierById, err := sh.Repo.GetbyId(id)

	if err != nil {
		fmt.Println(err)
	}

	jsonData, err := json.Marshal(supplierById)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}

func (sp *SupplierHandler) GetMenu(w http.ResponseWriter, r *http.Request) {

}
