package domain

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTreatmentGetters(t *testing.T) {
	date, err := time.Parse("2006-01-02", "2024-01-17")
	require.NoError(t, err)
	lastModified := date.AddDate(0, 0, 9)
	treatment := Treatment{
		Type:         "Medical Appointment",
		DateStart:    date,
		LastModified: lastModified,
	}

	expectedTreatmentName := fmt.Sprintf("Medical Appointment (2024-01-17)")
	assert.Equal(t, expectedTreatmentName, treatment.GetName())
	assert.Equal(t, lastModified, treatment.GetDate())
}

func TestTreatment_UnmarshalJSON(t *testing.T) {
	currentTime := time.Now().Truncate(0)
	owner := "Los Palmeras"
	// The following comments are from newest to oldest
	comment1 := Comment{
		DateAdded:   currentTime,
		Information: "Como hago compañero pa' decirle que no he podido olvidarla",
		Owner:       owner,
	}
	comment2 := Comment{
		DateAdded:   currentTime.Add(-1 * time.Hour),
		Information: "Que por más que lo intente sus recuerdos siempre habitan en mi mente",
		Owner:       owner,
	}
	comment3 := Comment{
		DateAdded:   currentTime.Add(-2 * time.Hour),
		Information: "Que no puedo pasar siquiera un día sin verla así sea desde lejos",
		Owner:       owner,
	}
	comment4 := Comment{
		DateAdded:   currentTime.Add(-3 * time.Hour),
		Information: "Que siento enloquecer al verla alegre, sonreír y no es conmigo",
		Owner:       owner,
	}
	comment5 := Comment{
		DateAdded:   currentTime.Add(-4 * time.Hour),
		Information: "Que siento enloquecer al verla alegre, sonreír y no es conmigo",
		Owner:       owner,
	}
	comment6 := Comment{
		DateAdded:   currentTime.Add(-5 * time.Hour),
		Information: "Olvídala",
		Owner:       owner,
	}
	comment7 := Comment{
		DateAdded:   currentTime.Add(-6 * time.Hour),
		Information: " No es fácil para mí, por eso quiero hablarle. Si es preciso rogarle que regrese a mi vida",
		Owner:       owner,
	}

	treatment := Treatment{
		ID:   "69",
		Type: "Medical Appointment",
		Comments: []Comment{
			comment7,
			comment5,
			comment4,
			comment2,
			comment1,
			comment3,
			comment6,
		},
		DateStart:    currentTime.Add(-6 * time.Hour),
		LastModified: currentTime,
		DateEnd:      nil,
		NextTurn:     nil,
	}

	// comments are sorted by DateAdded
	expectedTreatment := Treatment{
		ID:   "69",
		Type: "Medical Appointment",
		Comments: []Comment{
			comment1,
			comment2,
			comment3,
			comment4,
			comment5,
			comment6,
			comment7,
		},
		DateStart:    currentTime.Add(-6 * time.Hour),
		LastModified: currentTime,
		DateEnd:      nil,
		NextTurn:     nil,
	}

	rawTreatment, err := json.Marshal(treatment)
	require.NoError(t, err)

	var result Treatment
	err = json.Unmarshal(rawTreatment, &result)
	require.NoError(t, err)
	assert.Equal(t, expectedTreatment, result)
}

func TestComment_GetCommentMessage(t *testing.T) {
	date, err := time.Parse("2006-01-02", "2024-01-17")
	require.NoError(t, err)
	info := "Flasheaste amor"
	owner := "Agapornis ft Hernán y La Champions Liga"
	comment := Comment{
		DateAdded:   date,
		Information: info,
		Owner:       owner,
	}

	expectedCommentMessage := fmt.Sprintf("2024-01-17 by %s: %s", owner, info)
	assert.Equal(t, date, comment.GetDate())
	assert.Equal(t, expectedCommentMessage, comment.GetCommentMessage())
}
