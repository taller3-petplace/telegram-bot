package domain

import "time"

type Treatment struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Comments  []comment  `json:"comments"`
	DateStart time.Time  `json:"date_start"`
	DateEnd   *time.Time `json:"date_end"`
	NextTurn  *time.Time `json:"next_dose"`
}

type comment struct {
	DateAdded   time.Time `json:"date_added"`
	Information string    `json:"information"`
	Owner       string    `json:"owner"`
}

type Vaccine struct {
	Name      string     `json:"name"`
	FirstDose time.Time  `json:"first_dose"`
	LastDose  time.Time  `json:"last_dose"`
	NextDose  *time.Time `json:"next_dose"`
}
