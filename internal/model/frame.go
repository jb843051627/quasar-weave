package model

import "time"

type CalibrationFrame struct {
	ID            string     `json:"id"`
	ObservationID string     `json:"observation_id"`
	AntennaID     string     `json:"antenna_id"`
	Sequence      int        `json:"sequence"`
	CapturedAt    time.Time  `json:"captured_at"`
	Signal        float64    `json:"signal"`
	Drift         float64    `json:"drift"`
	Completeness  float64    `json:"completeness"`
	PayloadHash   string     `json:"payload_hash"`
	State         FrameState `json:"state"`
	ReceivedAt    time.Time  `json:"received_at"`
	CheckedAt     time.Time  `json:"checked_at"`
	RejectReason  string     `json:"reject_reason"`
}

type FrameInput struct {
	ID            string    `json:"id"`
	ObservationID string    `json:"observation_id"`
	AntennaID     string    `json:"antenna_id"`
	Sequence      int       `json:"sequence"`
	CapturedAt    time.Time `json:"captured_at"`
	Signal        float64   `json:"signal"`
	Drift         float64   `json:"drift"`
	Completeness  float64   `json:"completeness"`
	PayloadHash   string    `json:"payload_hash"`
}

func (f CalibrationFrame) ValidMeasurements() bool {
	return f.Signal >= 0 && f.Drift >= 0 && f.Completeness >= 0 && f.Completeness <= 1
}

func (f CalibrationFrame) Age(now time.Time) time.Duration {
	if f.CapturedAt.IsZero() {
		return 0
	}
	return now.Sub(f.CapturedAt)
}
