package model

type Supplier struct {
	ID           int
	ExtID        int          `json:"id"`
	Name         string       `json:"name"`
	Type         Type         `json:"type"`
	Image        string       `json:"image"`
	WorkingHours WorkingHours `json:"workingHours"`
}

type SupplierFromAPI struct {
	ID           int
	ExtID        int          `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Image        string       `json:"image"`
	WorkingHours WorkingHours `json:"workingHours"`
}

type Type struct {
	ID   int    `json:"id"`   // The unique identifier for the type
	Type string `json:"type"` // The name or category of the type
}

type WorkingHours struct {
	Opening string `json:"opening"`
	Closing string `json:"closing"`
}
type Response struct {
	Limit     int               `json:"limit"`
	Page      int               `json:"page"`
	Total     int               `json:"total"`
	Suppliers []SupplierFromAPI `json:"suppliers"`
}
