package domain

import (
	"encoding/json"
	"fmt"
	"telegram-bot/internal/utils"
	"time"
)

type Treatment struct {
	ID           string     `json:"id"`
	Type         string     `json:"type"`
	Comments     []Comment  `json:"comments"`
	DateStart    time.Time  `json:"date_start"`
	LastModified time.Time  `json:"last_modified"`
	DateEnd      *time.Time `json:"date_end"`
	NextTurn     *time.Time `json:"next_dose"`
}

func (t *Treatment) UnmarshalJSON(rawData []byte) error {
	var treatment struct {
		ID           string     `json:"id"`
		Type         string     `json:"type"`
		Comments     []Comment  `json:"comments"`
		DateStart    time.Time  `json:"date_start"`
		LastModified time.Time  `json:"last_modified"`
		DateEnd      *time.Time `json:"date_end"`
		NextTurn     *time.Time `json:"next_dose"`
	}
	err := json.Unmarshal(rawData, &treatment)
	if err != nil {
		return err
	}

	comments := treatment.Comments
	utils.SortElementsByDate(comments)

	t.ID = treatment.ID
	t.Type = treatment.Type
	t.Comments = comments
	t.DateStart = treatment.DateStart
	t.LastModified = treatment.LastModified
	t.DateEnd = treatment.DateEnd
	t.NextTurn = treatment.NextTurn

	return nil
}

// GetDate returns the date on which the treatment was modified by last time
func (t Treatment) GetDate() time.Time {
	return t.LastModified
}

// GetName the name of a treatment consists in the Type field follow by the DateStart, eg: medical appointment (2024-01-21)
func (t Treatment) GetName() string {
	return fmt.Sprintf("%s (%s)", t.Type, utils.DateToString(t.DateStart))
}

type Comment struct {
	DateAdded   time.Time `json:"date_added"`
	Information string    `json:"information"`
	Owner       string    `json:"owner"`
}

// GetCommentMessage returns a string with all the data about the comment. The format of the message is:
// DateAdded by Owner: Information. Eg: 2024/01/05 by McFly: some information
func (c Comment) GetCommentMessage() string {
	return fmt.Sprintf("%s by %s: %s", utils.DateToString(c.DateAdded), c.Owner, c.Information)
}

// GetDate returns the date on which the comment was added
func (c Comment) GetDate() time.Time {
	return c.DateAdded
}

type Vaccine struct {
	Name          string
	AmountOfDoses int
	FirstDose     time.Time
	LastDose      time.Time
}

// GetDate returns the date on which the last dose was applied
func (v Vaccine) GetDate() time.Time {
	return v.LastDose
}

// VaccineResponse response from Treatments service
type VaccineResponse struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

// GetDate returns the date on which the vaccine was applied
func (v VaccineResponse) GetDate() time.Time {
	return v.Date
}
