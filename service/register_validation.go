package service

import (
	"errors"
	"food_delivery/request"
	"unicode"
)

// validateInput validates the registration input fields without using regex
func ValidateInput(req *request.RegisterRequest) error {
	// Check for empty fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		return errors.New("all fields are required")
	}

	// Validate email (simple check for @ and .)
	if !IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate password strength (no regex)
	if !IsStrongPassword(req.Password) {
		return errors.New("password must be at least 8 characters long and contain at least one number, one uppercase letter, and one special character")
	}

	// Validate phone number (basic check for length)
	if len(req.Phone) < 10 || len(req.Phone) > 15 {
		return errors.New("phone number must be between 10 and 15 digits")
	}

	return nil
}

func IsValidEmail(email string) bool {
	// Basic email check: it must have one '@' and one '.' after the '@'
	at := false
	dot := false

	for i, char := range email {
		if char == '@' {
			at = true
		} else if at && char == '.' && i > 0 {
			dot = true
		}
	}

	return at && dot
}

// isStrongPassword checks for password strength (no regex)
func IsStrongPassword(password string) bool {
	var hasMinLen bool = len(password) >= 8
	var hasUppercase, hasLowercase, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Return true if all conditions are met
	return hasMinLen && hasUppercase && hasLowercase && hasNumber && hasSpecial
}
