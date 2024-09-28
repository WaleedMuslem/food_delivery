package repository

import (
	"database/sql"
	"food_delivery/model"
)

type AddressRepository struct {
	Db *sql.DB
}

type IAddressRepository interface {
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return AddressRepository{Db: db}
}

func (ar AddressRepository) Create(address *model.Address, userId uint) (int, error) {

	var addressId int
	// Prepare the SQL statement
	stmt, err := ar.Db.Prepare(`INSERT INTO addresses (user_id, floor, apartment, street, city, zip_code) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING address_id`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the statement
	err = stmt.QueryRow(userId, address.Floor, address.Apartment, address.Street, address.City, address.Zip).Scan(&addressId)
	if err != nil {
		return 0, err
	}

	return addressId, nil

}

func (ar AddressRepository) Get(userId uint) (model.Address, error) {

	query := `SELECT street, city, zip_code, floor, apartment
	FROM addresses
	WHERE user_id = $1
	ORDER BY address_id DESC
	LIMIT 1;`

	var address model.Address

	err := ar.Db.QueryRow(query, userId).Scan(&address.Street, &address.City, &address.Zip, &address.Floor, &address.Apartment)
	if err != nil {
		return model.Address{}, err
	}

	return address, nil

}
