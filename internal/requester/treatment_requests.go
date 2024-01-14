package requester

import "telegram-bot/internal/domain"

// GetTreatmentsByPetID fetches the last 5 treatments of the given pet.
// The treatments go from new ones to old ones
func (r *Requester) GetTreatmentsByPetID(petID int) ([]domain.Treatment, error) {
	return nil, nil
}

// GetTreatment fetches all the information about the given treatment
func (r *Requester) GetTreatment(treatmentID int) ([]domain.Treatment, error) {
	return nil, nil
}

// GetVaccines fetches all the vaccines that were applied to the pet
func (r *Requester) GetVaccines(petID int) ([]domain.Vaccine, error) {
	return nil, nil
}
