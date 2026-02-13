package status

import (
	"testing"
	datetime "time"

	"github.com/stretchr/testify/require"
)

func TestTime_UnmarshalJSON_ValidDate(t *testing.T) {
	rawTime := "\"2014-05-03T01:22:07.286Z\""
	time := new(Time)
	err := time.UnmarshalJSON([]byte(rawTime))
	require.NoError(t, err)

	require.Equal(t, 2014, time.Year())
	require.Equal(t, datetime.May, time.Month())
	require.Equal(t, 3, time.Day())
	require.Equal(t, 1, time.Hour())
	require.Equal(t, 22, time.Minute())
	require.Equal(t, 7, time.Second())
}

func TestTime_UnmarshalJSON_Null(t *testing.T) {
	rawTime := "\"null\""
	time := new(Time)
	err := time.UnmarshalJSON([]byte(rawTime))
	require.NoError(t, err)
	require.Nil(t, time.Time)
}

func TestTime_UnmarshalJSON_Invalid(t *testing.T) {
	rawTime := "\"not a time\""
	time := new(Time)
	err := time.UnmarshalJSON([]byte(rawTime))
	require.Error(t, err)
	require.Nil(t, time.Time)
}

func TestTime_MarshalJSON_Time(t *testing.T) {
	rawTime := "\"2014-05-03T01:22:07.286Z\""
	time := new(Time)
	err := time.UnmarshalJSON([]byte(rawTime))
	require.NoError(t, err)
	timeString, err := time.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "\"2014-05-03 1:22 AM\"", string(timeString))
}

func TestTime_MarshalJSON_Null(t *testing.T) {
	time := new(Time)
	timeString, err := time.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "null", string(timeString))
}
func TestIncidentUpdate_UnmarshalJSON_StripsBrTags(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{
			name:     "br with space and slash",
			json:     `{"body":"some text<br />more text","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "some text more text",
		},
		{
			name:     "self-closing br",
			json:     `{"body":"some text<br/>more text","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "some text more text",
		},
		{
			name:     "bare br",
			json:     `{"body":"some text<br>more text","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "some text more text",
		},
		{
			name:     "multiple br tags",
			json:     `{"body":"first<br />second<br/>third<br>fourth","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "first second third fourth",
		},
		{
			name:     "no br tags",
			json:     `{"body":"just plain text","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "just plain text",
		},
		{
			name:     "trailing br tag",
			json:     `{"body":"message text<br />","status":"investigating","created_at":"2014-05-03T01:22:07.286Z"}`,
			expected: "message text ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var update IncidentUpdate
			err := update.UnmarshalJSON([]byte(tt.json))
			require.NoError(t, err)
			require.Equal(t, tt.expected, update.Update)
		})
	}
}

func TestTime_String_Null(t *testing.T) {
	time := new(Time)
	require.Equal(t, "", time.String())
}
