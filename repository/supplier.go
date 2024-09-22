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
	suppliers := []model.Supplier{}

	// Updated query to join with the type table
	result, err := sr.Db.Query(`
		SELECT s.id, s.name, s.image, s.opening, s.closing, s.ext_id, s.type_id, t.type 
		FROM suppliers s 
		JOIN supplier_type t ON s.type_id = t.id 
		ORDER BY s.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		supplier := model.Supplier{}
		supplierType := model.Type{} // Assuming there's a Type struct in model

		// Scanning both supplier and type data
		err := result.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Image,
			&supplier.WorkingHours.Opening,
			&supplier.WorkingHours.Closing,
			&supplier.ExtID,
			&supplierType.ID,   // Scanning type_id into Type struct's ID
			&supplierType.Type, // Scanning the type field
		)
		if err != nil {
			return nil, err
		}

		// Assign the scanned type to the supplier's Type field
		supplier.Type = supplierType

		suppliers = append(suppliers, supplier)
	}

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
