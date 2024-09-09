package handler

import (
	"encoding/json"
	"fmt"
	"food_delivery/repository"
	"net/http"
	"strconv"
)

type MenuHandler struct {
	Repo *repository.MenuRepository
}

func NewMenuController(repo repository.MenuRepository) MenuHandler {
	return MenuHandler{Repo: &repo}
}

// func (mh *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var menu model.Product

// 	err = json.Unmarshal(body, &menu)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = mh.Repo.Create(menu)
// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonData, err := json.Marshal(menu)
// 	if err != nil {
// 		fmt.Println("Failed to marshal JSON ")
// 	}

// 	w.WriteHeader(200)
// 	w.Write(jsonData)
// }

func (mh *MenuHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	// fmt.Println(idString)
	id, _ := strconv.Atoi(idString)

	menus, err := mh.Repo.GetAll(id)
	if err != nil {
		fmt.Println(err)
	}

	jsonData, err := json.Marshal(menus)
	if err != nil {
		fmt.Println("Failed to marshal JSON ")
	}

	w.WriteHeader(200)
	w.Write(jsonData)
}
