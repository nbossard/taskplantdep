package model

import (
	"strings"
)

// ExportedTask is the model for a task that is exported by taskwarrior in JSON.
type ExportedTask struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Due         string   `json:"due"`
	End         string   `json:"end"`
	Entry       string   `json:"entry"`
	Modified    string   `json:"modified"`
	Project     string   `json:"project"`
	Status      string   `json:"status"`
	UUID        string   `json:"uuid"`
	Wait        string   `json:"wait"`
	Depends     []string `json:"depends"`
	Urgency     float64  `json:"urgency"`
	Tags        []string `json:"tags"`
}

// GetUUIDCleaned returns the UUID of the task, but without the "-" characters.
// To be compatible with PlantUML objects.
func (t *ExportedTask) GetUUIDCleaned() string {
	return MakeOneUUIDCompatible(t.UUID)
}

// MakeOneUUIDCompatible transforms provided UUID to be compatible with PlantUML objects.
// Removes all - from UUIDs.
func MakeOneUUIDCompatible(parUUID string) string {
	return strings.Replace(parUUID, "-", "", -1)
}

// GetColorUrgency returns the color according to the urgency of the task.
// See list of colors here: https://plantuml.com/color
func (t *ExportedTask) GetColorUrgency() string {
	if t.Urgency < 1 {
		return "white"
	}
	if t.Urgency < 5 {
		return "lightGreen"
	}
	if t.Urgency < 10 {
		return "gold"
	}
	if t.Urgency < 15 {
		return "orange"
	}
	return "red"
}
