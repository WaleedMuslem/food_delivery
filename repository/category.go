package repository

import (
	"database/sql"
	"food_delivery/model"
)

type ICategory interface {
	GetAll() ([]model.Category, error)
}

type CategoryRepository struct {
	Db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return CategoryRepository{Db: db}
}

func (cr CategoryRepository) GetAll() ([]model.Category, error) {
	categories := []model.Category{}

	// Query to get categories with the count of products
	query := `
		SELECT c.category_id, c.category_name, c.image, COUNT(p.id) AS product_count
		FROM category c
		LEFT JOIN products p ON c.category_id = p.category_id
		GROUP BY c.category_id, c.category_name, c.image
	`
	result, err := cr.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close() // Ensure the result set is closed after processing

	for result.Next() {
		category := model.Category{}

		err := result.Scan(&category.Category_id, &category.Category_name, &category.Category_image, &category.ProductCount)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	result.Close()

	return categories, nil
}
