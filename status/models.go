package status

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type SystemStatus struct {
	Components []Components `json:"components"`
	Incidents  []Incidents  `json:"incidents"`
}

type Components struct {
	ID        string          `json:"id"`
	Component string          `json:"name"`
	Status    ComponentStatus `json:"status"`
}

type Incidents struct {
	Status          IncidentStatus   `json:"status"`
	ID              string           `json:"id"`
	IncidentUpdates []IncidentUpdate `json:"incident_updates"`
}

type IncidentUpdate struct {
	Status    IncidentStatus `json:"status"`
	Update    string         `json:"body"`
	Timestamp *Time          `json:"created_at"`
}

type ComponentStatus string
type IncidentStatus string

var brTagRegex = regexp.MustCompile(`<br\s*/?>`)

func (iu *IncidentUpdate) UnmarshalJSON(b []byte) error {
	type Alias IncidentUpdate
	var alias Alias
	if err := json.Unmarshal(b, &alias); err != nil {
		return err
	}
	alias.Update = brTagRegex.ReplaceAllString(alias.Update, " ")
	*iu = IncidentUpdate(alias)
	return nil
}

type Time struct {
	*time.Time
}

func (t *Time) String() string {
	if t.Time == nil {
		return ""
	}
	return t.Time.Format("2006-01-02 3:04 PM")
}

func (t *Time) UnmarshalJSON(b []byte) error {
	timeAsString := string(b)
	timeAsString = strings.Trim(timeAsString, "\"")
	if timeAsString == "null" {
		t.Time = nil
		return nil
	}
	tim, err := time.Parse("2006-01-02T15:04:05.999Z", timeAsString)
	if err != nil {
		return err
	}
	t.Time = &tim
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", t.String())), nil
}
