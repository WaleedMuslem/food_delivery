package model

type Address struct {
	Floor     string `json:"floor"`
	Apartment string `json:"apartment"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Zip       string `json:"zip"`
}
