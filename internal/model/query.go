package model

import "time"

type FrameFilter struct {
	ObservationID string
	AntennaID     string
	State         FrameState
	Since         time.Time
	Limit         int
}

type AlertFilter struct {
	ObservationID string
	State         AlertState
	Severity      string
	Limit         int
}

type NoteInput struct {
	ObservationID string `json:"observation_id"`
	AlertID       string `json:"alert_id"`
	Operator      string `json:"operator"`
	Body          string `json:"body"`
}

type GateInput struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	MinSignal       float64 `json:"min_signal"`
	MaxDrift        float64 `json:"max_drift"`
	MinCompleteness float64 `json:"min_completeness"`
	WindowSeconds   int     `json:"window_seconds"`
	Active          bool    `json:"active"`
}
