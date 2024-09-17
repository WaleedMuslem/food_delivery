package model

type Supplier struct {
	ID           int
	ExtID        int          `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Image        string       `json:"image"`
	WorkingHours WorkingHours `json:"workingHours"`
}

type WorkingHours struct {
	Opening string `json:"opening"`
	Closing string `json:"closing"`
}
type Response struct {
	Limit     int        `json:"limit"`
	Page      int        `json:"page"`
	Total     int        `json:"total"`
	Suppliers []Supplier `json:"suppliers"`
}
