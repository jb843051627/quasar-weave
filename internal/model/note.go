package model

import "time"

type OperatorNote struct {
	ID            string    `json:"id"`
	ObservationID string    `json:"observation_id"`
	AlertID       string    `json:"alert_id"`
	Operator      string    `json:"operator"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"created_at"`
}

type Event struct {
	ID        string    `json:"id"`
	Subject   string    `json:"subject"`
	Action    string    `json:"action"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type HealthSummary struct {
	AntennaCount       int     `json:"antenna_count"`
	HealthyAntennas    int     `json:"healthy_antennas"`
	ActiveObservations int     `json:"active_observations"`
	OpenAlerts         int     `json:"open_alerts"`
	QueuedFrames       int     `json:"queued_frames"`
	AverageQuality     float64 `json:"average_quality"`
}
