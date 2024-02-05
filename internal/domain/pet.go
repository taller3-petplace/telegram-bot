package domain

import "time"

// PetRequest request body to create a pet register
type PetRequest struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	RegisterDate time.Time `json:"register_date"`
	BirthDate    string    `json:"birth_date"`
	OwnerID      string    `json:"owner_id"`
}

// PetDataIdentifier brief data to identify a pet
type PetDataIdentifier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// PetData general data for a pet. Does not contain anything about treatments
type PetData struct {
	PetDataIdentifier
	BirthDate time.Time `json:"birth_date"`
	Race      string    `json:"race,omitempty"`
}

// PetsResponse groups data from different pets for a given user
type PetsResponse struct {
	PetsData []PetData `json:"results"`
	Paging   Paging    `json:"paging"`
}

type Paging struct {
	Total  uint `json:"total"`
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}
