package bot

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var missingDotsInTypeField = `
Name:
Birth Date:
Type
`

var fieldsInOtherOrder = `
Birth Date: 2023/12/10

Name: Ringo

Type: DOG
`

var validEmptyForm = `
Name:
Birth Date:
Type:
`

var validFormWithMultipleNewlines = `
Name: Pumba



Birth Date: 2023/12/10




Type: dog
`

var validFormWithoutNewlines = "Name: Plumitas Birth Date: 10/12/2023 Type: DUCK"

var validNormalForm = `
Name: Ringo

Birth Date: 2023/12/10

Type: DOG
`

func TestExtractPetDataErrorDueToInvalidForm(t *testing.T) {
	fieldsTags := []string{nameTag, birthDateTag, typeTag}
	testCases := []struct {
		Name          string
		Form          string
		FieldTags     []string
		ExpectedError error
	}{
		{
			Name:          "Empty form",
			Form:          missingDotsInTypeField,
			FieldTags:     fieldsTags,
			ExpectedError: errInvalidForm,
		},
		{
			Name:          "Form with fields in other order",
			Form:          fieldsInOtherOrder,
			FieldTags:     fieldsTags,
			ExpectedError: errInvalidForm,
		},
		{
			Name:          "Missing field tags",
			Form:          validNormalForm,
			FieldTags:     []string{nameTag, birthDateTag, typeTag, "random-tag"},
			ExpectedError: errMissingFormField,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			petData, err := extractPetData(testCase.Form, testCase.FieldTags...)
			assert.Nil(t, petData)
			assert.ErrorIs(t, err, testCase.ExpectedError)
		})
	}
}

func TestExtractPetDataCorrectly(t *testing.T) {
	testCases := []struct {
		Name              string
		Form              string
		ExpectedName      string
		ExpectedBirthDate string
		ExpectedType      string
	}{
		{
			Name:              "Empty form",
			Form:              validEmptyForm,
			ExpectedName:      "",
			ExpectedBirthDate: "",
			ExpectedType:      "",
		},
		{
			Name:              "Valid form with multiple newlines",
			Form:              validFormWithMultipleNewlines,
			ExpectedName:      "Pumba",
			ExpectedBirthDate: "2023/12/10",
			ExpectedType:      "dog",
		},
		{
			Name:              "Valid form without newlines between fields",
			Form:              validFormWithoutNewlines,
			ExpectedName:      "Plumitas",
			ExpectedBirthDate: "10/12/2023",
			ExpectedType:      "DUCK",
		},
		{
			Name:              "Valid normal form",
			Form:              validNormalForm,
			ExpectedName:      "Ringo",
			ExpectedBirthDate: "2023/12/10",
			ExpectedType:      "DOG",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			petData, err := extractPetData(testCase.Form, nameTag, birthDateTag, typeTag)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedName, petData[nameTag])
			assert.Equal(t, testCase.ExpectedBirthDate, petData[birthDateTag])
			assert.Equal(t, testCase.ExpectedType, petData[typeTag])
		})
	}
}
