package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"food_delivery/model"
	"food_delivery/request"
	"food_delivery/service"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db *sql.DB
}

type IUserRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(id uint) (*model.User, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{Db: db}
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {

	var userById model.User
	err := ur.Db.QueryRow("SELECT * FROM users WHERE email = $1",
		&email).Scan(&userById.ID, &userById.FirstName, &userById.LastName, &userById.Email, &userById.Password, &userById.Phone)
	if err != nil {
		return nil, err
	}

	return &userById, err
}

func (ur *UserRepository) GetUserById(id uint) (*model.User, error) {

	var userById model.User
	err := ur.Db.QueryRow("SELECT * FROM users WHERE user_id = $1",
		&id).Scan(&userById.ID, &userById.FirstName, &userById.LastName, &userById.Email, &userById.Password, &userById.Phone)
	if err != nil {
		return nil, err
	}

	return &userById, err
}

func (ur *UserRepository) RegisterUser(req *request.RegisterRequest) error {
	if err := service.ValidateInput(req); err != nil {
		return err
	}

	// Check if the username already exists
	var existingUser model.User
	err := ur.Db.QueryRow("SELECT email FROM users WHERE email = $1", req.Email).
		Scan(&existingUser.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking username: %v", err)
	}
	if existingUser.Email != "" {
		return errors.New("username already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = ur.Db.Exec("INSERT INTO users (first_name, last_name, email, password, phone) VALUES ($1, $2, $3, $4, $5)", &req.FirstName, &req.LastName, &req.Email, string(hashedPassword), req.Phone)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	return nil
}

func (ur *UserRepository) Logout(userId uint) error {

	_, err := ur.Db.Exec("UPDATE jwt_tokens SET revoked = true WHERE user_id = $1", userId)

	return err
}
