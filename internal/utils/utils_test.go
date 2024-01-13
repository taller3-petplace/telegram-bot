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
