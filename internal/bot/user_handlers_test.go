package bot

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var missingDotsInHourField = `
Hour
End Date:
`

var alarmFieldsInOtherOrder = `
End Date: 2023/12/10

Hour: 10:30
`

var validEmptyAlarmForm = `
Hour:
End Date:
`

var validAlarmFormWithMultipleNewlines = `
Hour: 10:00



End Date: 2023/12/10
`

var validAlarmFormWithoutNewlines = "Hour: 9:30 End Date: 10/12/2023"

var validNormalAlarmForm = `
Hour: 1:00

End Date: 2023/12/10
`

var validAlarmFormWithMultipleSpacesBeforeFieldValues = `
Hour:      1:00

End Date:     2023/12/10
`

var validAlarmFormWithNotApplicableEndDate = `
Hour: 2:00

End Date: N/A
`

func TestExtractAlarmErrorDueToInvalidForm(t *testing.T) {
	fieldsTags := []string{hourTag, endDateTag}
	testCases := []struct {
		Name          string
		Form          string
		FieldTags     []string
		ExpectedError error
	}{
		{
			Name:          "Empty form",
			Form:          missingDotsInHourField,
			FieldTags:     fieldsTags,
			ExpectedError: errInvalidForm,
		},
		{
			Name:          "Form with fields in other order",
			Form:          alarmFieldsInOtherOrder,
			FieldTags:     fieldsTags,
			ExpectedError: errInvalidForm,
		},
		{
			Name:          "Missing field tags",
			Form:          validNormalAlarmForm,
			FieldTags:     []string{hourTag, endDateTag, "random-tag"},
			ExpectedError: errMissingFormField,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			alarmData, err := extractAlarmData(testCase.Form, testCase.FieldTags...)
			assert.Nil(t, alarmData)
			assert.ErrorIs(t, err, testCase.ExpectedError)
		})
	}
}

func TestExtractAlarmDataCorrectly(t *testing.T) {
	testCases := []struct {
		Name            string
		Form            string
		ExpectedHour    string
		ExpectedEndDate string
	}{
		{
			Name:            "Empty form",
			Form:            validEmptyAlarmForm,
			ExpectedHour:    "",
			ExpectedEndDate: "",
		},
		{
			Name:            "Valid form with multiple newlines",
			Form:            validAlarmFormWithMultipleNewlines,
			ExpectedHour:    "10:00",
			ExpectedEndDate: "2023/12/10",
		},
		{
			Name:            "Valid form without newlines between fields",
			Form:            validAlarmFormWithoutNewlines,
			ExpectedHour:    "9:30",
			ExpectedEndDate: "10/12/2023",
		},
		{
			Name:            "Valid normal form",
			Form:            validNormalAlarmForm,
			ExpectedHour:    "1:00",
			ExpectedEndDate: "2023/12/10",
		},
		{
			Name:            "Valid normal form with multiple spaces before field values",
			Form:            validAlarmFormWithMultipleSpacesBeforeFieldValues,
			ExpectedHour:    "1:00",
			ExpectedEndDate: "2023/12/10",
		},
		{
			Name:            "End date with not applicable value",
			Form:            validAlarmFormWithNotApplicableEndDate,
			ExpectedHour:    "2:00",
			ExpectedEndDate: notApplicable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			alarmData, err := extractAlarmData(testCase.Form, hourTag, endDateTag)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedHour, alarmData[hourTag])
			assert.Equal(t, testCase.ExpectedEndDate, alarmData[endDateTag])
		})
	}
}
