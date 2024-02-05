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
Message:
Hours:
Start Date:
End Date:
`

var validAlarmFormWithMultipleNewlines = `
Message: Olvídala. No es fácil para mi, por eso quiero hablarle, si es preciso rogarle que regrese a mi vida



Hours: 10:00



Start Date: 2023/12/10



End Date: 2023/12/10
`

var validAlarmFormWithoutNewlines = "Message: Hola perdida Hours: 9:30, 12:30 Start Date: 10/12/2023 End Date: 10/12/2023"

var validNormalAlarmForm = `
Message: hola que tal tu como estas? dime si eres feliz

Hours: 10:30, 22:30

Start Date: 2023/12/10

End Date: 2023/12/10
`

var validAlarmFormWithMultipleSpacesBeforeFieldValues = `
Message:    si te invito una copa y me acerco a tu boca
Hours:      1:00

Start Date:     2023/12/10

End Date:     2023/12/10
`

var validAlarmFormWithNotApplicableEndDate = `
Message: Te pido de rodillas luna no te vayas

Hours: 2:00

Start Date: 2024/02/04

End Date: N/A
`

func TestExtractAlarmErrorDueToInvalidForm(t *testing.T) {
	fieldsTags := []string{hoursTag, endDateTag}
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
			FieldTags:     []string{hoursTag, endDateTag, "random-tag"},
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
		Name              string
		Form              string
		ExpectedMessage   string
		ExpectedHours     string
		ExpectedStartDate string
		ExpectedEndDate   string
	}{
		{
			Name:            "Empty form",
			Form:            validEmptyAlarmForm,
			ExpectedHours:   "",
			ExpectedEndDate: "",
		},
		{
			Name:              "Valid form with multiple newlines",
			Form:              validAlarmFormWithMultipleNewlines,
			ExpectedMessage:   "Olvídala. No es fácil para mi, por eso quiero hablarle, si es preciso rogarle que regrese a mi vida",
			ExpectedHours:     "10:00",
			ExpectedStartDate: "2023/12/10",
			ExpectedEndDate:   "2023/12/10",
		},
		{
			Name:              "Valid form without newlines between fields",
			Form:              validAlarmFormWithoutNewlines,
			ExpectedMessage:   "Hola perdida",
			ExpectedHours:     "9:30, 12:30",
			ExpectedStartDate: "10/12/2023",
			ExpectedEndDate:   "10/12/2023",
		},
		{
			Name:              "Valid normal form",
			Form:              validNormalAlarmForm,
			ExpectedMessage:   "hola que tal tu como estas? dime si eres feliz",
			ExpectedHours:     "10:30, 22:30",
			ExpectedStartDate: "2023/12/10",
			ExpectedEndDate:   "2023/12/10",
		},
		{
			Name:              "Valid normal form with multiple spaces before field values",
			Form:              validAlarmFormWithMultipleSpacesBeforeFieldValues,
			ExpectedMessage:   "si te invito una copa y me acerco a tu boca",
			ExpectedHours:     "1:00",
			ExpectedStartDate: "2023/12/10",
			ExpectedEndDate:   "2023/12/10",
		},
		{
			Name:              "End date with not applicable value",
			Form:              validAlarmFormWithNotApplicableEndDate,
			ExpectedMessage:   "Te pido de rodillas luna no te vayas",
			ExpectedHours:     "2:00",
			ExpectedStartDate: "2024/02/04",
			ExpectedEndDate:   notApplicable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			alarmData, err := extractAlarmData(testCase.Form, messageTag, hoursTag, startDateTag, endDateTag)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedMessage, alarmData[messageTag])
			assert.Equal(t, testCase.ExpectedHours, alarmData[hoursTag])
			assert.Equal(t, testCase.ExpectedStartDate, alarmData[startDateTag])
			assert.Equal(t, testCase.ExpectedEndDate, alarmData[endDateTag])
		})
	}
}
