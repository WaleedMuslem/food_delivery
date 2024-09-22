package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/repository"
	"net/http"
)

type CategoryHandler struct {
	Repo *repository.CategoryRepository
}

func NewcategoryController(repo repository.CategoryRepository) CategoryHandler {
	return CategoryHandler{Repo: &repo}
}

func (ch *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	categories, err := ch.Repo.GetAll()
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(categories)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	// p1, _ := bcrypt.GenerateFromPassword([]byte("123456789"), bcrypt.MinCost)
	// fmt.Println(string(p1))

	w.WriteHeader(200)
	w.Write(jsonData)
}
