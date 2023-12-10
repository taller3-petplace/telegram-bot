package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateDateType(t *testing.T) {
	testCases := []struct {
		Name         string
		Date         string
		ExpectsError bool
	}{
		{
			Name:         "Invalid format: day/month/year",
			Date:         "10/12/2023",
			ExpectsError: true,
		},
		{
			Name:         "Invalid format: month/day/year",
			Date:         "12/10/2023",
			ExpectsError: true,
		},
		{
			Name:         "Invalid format: year-month-year",
			Date:         "2023-12-10",
			ExpectsError: true,
		},
		{
			Name:         "Valid format",
			Date:         "2023/12/10",
			ExpectsError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := ValidateDateType(testCase.Date)
			if testCase.ExpectsError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidatePetType(t *testing.T) {
	testCases := []struct {
		Name         string
		PetType      string
		ExpectsError bool
	}{
		{
			Name:         "Invalid pet",
			PetType:      "Bad Bunny",
			ExpectsError: true,
		},
		{
			Name:         "Valid pet",
			PetType:      "rabbit",
			ExpectsError: false,
		},
		{
			Name:         "Valid pet in upper case",
			PetType:      "RABBIT",
			ExpectsError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := ValidatePetType(testCase.PetType)
			if testCase.ExpectsError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
