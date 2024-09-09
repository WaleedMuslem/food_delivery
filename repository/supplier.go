package repository

import (
	"database/sql"
	"food_delivery/model"

	_ "github.com/lib/pq"
)

type ISupplier interface {
	Create(supplier model.Supplier) error
	GetAll() ([]model.Supplier, error)
	GetbyId() (model.Supplier, error)
}

type SupplierRepository struct {
	Db *sql.DB
}

func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return SupplierRepository{Db: db}
}

func (sr SupplierRepository) Create(supplier model.Supplier) error {

	_, err := sr.Db.Exec(
		"INSERT INTO suppliers (id, name, type, image, opening, closing) values ($1, $2, $3, $4, $5, $6)",
		supplier.ID, supplier.Name, supplier.Type, supplier.Image, supplier.WorkingHours.Opening, supplier.WorkingHours.Closing,
	)

	return err
}

func (sr SupplierRepository) GetAll() ([]model.Supplier, error) {

	// fetchedSuppliers, err := service.FetchAllSuppliers()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	// // fmt.Println(fetchedSuppliers)
	// for _, supplier := range fetchedSuppliers {
	// 	sr.Create(supplier)
	// }

	suppliers := []model.Supplier{}

	result, err := sr.Db.Query("SELECT id, name, type, image, opening, closing FROM suppliers ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		supplier := model.Supplier{}

		err := result.Scan(&supplier.ID, &supplier.Name, &supplier.Type, &supplier.Image, &supplier.WorkingHours.Opening, &supplier.WorkingHours.Closing)
		if err != nil {
			return nil, err
		}

		suppliers = append(suppliers, supplier)
	}

	result.Close()

	return suppliers, nil
}

func (sr SupplierRepository) GetbyId(id int) (*model.Supplier, error) {

	var supplierById model.Supplier
	err := sr.Db.QueryRow("SELECT * FROM suppliers WHERE id = $1",
		&id).Scan(&supplierById.ID, &supplierById.Name, &supplierById.Type, &supplierById.Image, &supplierById.WorkingHours.Opening, &supplierById.WorkingHours.Closing)
	if err != nil {
		return nil, err
	}

	return &supplierById, err
}
