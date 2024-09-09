package repository

import (
	"database/sql"
	"food_delivery/model"
)

// type Product struct {
// 	ID          int
// 	Name        string
// 	Price       float64
// 	SupplierID  int
// 	CategoryID  int
// 	Image       string
// 	Ingredients []string
// }

type IMenu interface {
	GetAll() ([]model.Menu, error)
}

type MenuRepository struct {
	Db *sql.DB
}

func NewMenuRepository(db *sql.DB) MenuRepository {
	return MenuRepository{Db: db}
}

func (mr MenuRepository) GetAll(supplier_id int) ([]model.Menu, error) {

	menus := []model.Menu{}

	result, err := mr.Db.Query("SELECT id, name, price,supplier_id, image, type FROM products WHERE supplier_id = $1", supplier_id)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		menu := model.Menu{}

		err := result.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.SupplierID, &menu.Image, &menu.Type)
		if err != nil {
			return nil, err
		}

		menus = append(menus, menu)
	}

	result.Close()

	return menus, nil
}
