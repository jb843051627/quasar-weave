package model

import "time"

type AuditEntry struct {
	ID         string    `json:"id"`
	Subject    string    `json:"subject"`
	Action     string    `json:"action"`
	Actor      string    `json:"actor"`
	Before     string    `json:"before"`
	After      string    `json:"after"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Timeline struct {
	ObservationID string            `json:"observation_id"`
	Events        []AuditEntry      `json:"events"`
	Commands      []OperatorCommand `json:"commands"`
}

func (e AuditEntry) Valid() bool {
	return e.ID != "" && e.Subject != "" && e.Action != "" && !e.OccurredAt.IsZero()
}

func (t Timeline) LastAction() string {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[0].Action
}
