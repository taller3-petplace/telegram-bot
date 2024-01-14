package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	t.Run("Contains int works correctly", func(t *testing.T) {
		elements := []int{1, 2, 3, 4, 5}
		assert.True(t, Contains(elements, 3))
		assert.False(t, Contains(elements, 69))
	})

	t.Run("Contains string works correctly", func(t *testing.T) {
		elements := []string{"hola", "que tal?", "tu como estas", "dime si eres feliz"}
		assert.True(t, Contains(elements, "hola"))
		assert.False(t, Contains(elements, "feliz"))
	})

	t.Run("Contains float works correctly", func(t *testing.T) {
		elements := []float32{1.2, 69.8, 28.5}
		assert.True(t, Contains(elements, 69.8))
		assert.False(t, Contains(elements, 169.0))
	})
}

func TestCalculateYearsBetweenDates(t *testing.T) {
	pastTime := time.Now().AddDate(-3, 0, 0)
	futureTime := time.Now().AddDate(69, 0, 0)
	assert.Equal(t, 3, CalculateYearsBetweenDates(pastTime))
	assert.Equal(t, -69, CalculateYearsBetweenDates(futureTime))
}

type testElement struct {
	date time.Time
}

func (te testElement) GetDate() time.Time {
	return te.date
}

func TestSortElementsByDate(t *testing.T) {
	// The following definitions are ordered from oldest to most recent
	currentTime := time.Now()

	elem1 := currentTime.AddDate(-69, 0, 0)
	elem2 := currentTime.AddDate(-50, 0, 0)
	elem3 := currentTime.AddDate(-28, 0, 0)
	elem4 := currentTime

	elementsToSort := []testElement{
		{date: elem1},
		{date: elem2},
		{date: elem3},
		{date: elem4},
	}

	expectedOrder := []testElement{
		{date: elem4},
		{date: elem3},
		{date: elem2},
		{date: elem1},
	}

	SortElementsByDate(elementsToSort)
	assert.Equal(t, expectedOrder, elementsToSort)
}

func TestDateToString(t *testing.T) {
	date := time.Date(2001, 9, 11, 0, 0, 0, 0, time.UTC)
	expectedString := "2001-09-11"
	assert.Equal(t, expectedString, DateToString(date))
}
