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
func TestTime_String_Null(t *testing.T) {
	time := new(Time)
	require.Equal(t, "", time.String())
}
