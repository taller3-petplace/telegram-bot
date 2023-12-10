package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserInfo(t *testing.T) {
	testUserID := int64(69)
	testCases := []struct {
		Name         string
		UserName     string
		ExpectedName string
	}{
		{
			Name:         "If name is empty, the name of the user will be 'unknown'",
			UserName:     "",
			ExpectedName: "Unknown",
		},
		{
			Name:         "Name is set correctly",
			UserName:     "Paul Vazo",
			ExpectedName: "Paul Vazo",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			userInfo := NewUserInfo(testUserID, testCase.UserName)
			assert.Equal(t, testUserID, userInfo.GetUserID())
			assert.Equal(t, testCase.ExpectedName, userInfo.GetName())
		})
	}
}
