package template

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// TestPetForm checks that the structure of the form does not change. It's really dummy
func TestPetForm(t *testing.T) {
	regex := regexp.MustCompile(`Name:[^\n]+\n+Birth Date:[^\n]+\n+Type:.*`)
	petForm := RegisterPet()
	assert.True(t, regex.MatchString(petForm))
}
