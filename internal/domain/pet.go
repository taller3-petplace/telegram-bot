package domain

import "time"

// PetRequest request body to create a pet register
type PetRequest struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	RegisterDate time.Time `json:"register_date"`
	BirthDate    string    `json:"birth_date"`
	OwnerID      int64     `json:"owner_id"`
}

// PetDataSummary contains brief information about the pet
type PetDataSummary struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birth_date"`
	Type      string    `json:"type"`
}

type PetsData struct {
	PetsData []PetDataSummary `json:"results"`
	Paging   Paging           `json:"paging"`
}

type Paging struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}
