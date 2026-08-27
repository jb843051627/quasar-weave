package model

import "time"

type Observation struct {
	ID             string            `json:"id"`
	Target         string            `json:"target"`
	RequestedBy    string            `json:"requested_by"`
	Status         ObservationStatus `json:"status"`
	StartAt        time.Time         `json:"start_at"`
	EndAt          time.Time         `json:"end_at"`
	ExpectedFrames int               `json:"expected_frames"`
	ReceivedFrames int               `json:"received_frames"`
	QualityScore   float64           `json:"quality_score"`
	GateID         string            `json:"gate_id"`
	FailureReason  string            `json:"failure_reason"`
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type ObservationInput struct {
	ID             string `json:"id"`
	Target         string `json:"target"`
	RequestedBy    string `json:"requested_by"`
	ExpectedFrames int    `json:"expected_frames"`
	GateID         string `json:"gate_id"`
}

type ObservationFilter struct {
	Status ObservationStatus
	Target string
	Limit  int
}

func (o Observation) Complete() bool {
	return o.Status == ObservationQualified || o.Status == ObservationFailed || o.Status == ObservationArchived
}

func (o Observation) Progress() float64 {
	if o.ExpectedFrames <= 0 {
		return 0
	}
	value := float64(o.ReceivedFrames) / float64(o.ExpectedFrames)
	if value > 1 {
		return 1
	}
	return value
}

func (o *Observation) RecordFrame(score float64, now time.Time) {
	o.ReceivedFrames++
	if o.ReceivedFrames == 1 || score < o.QualityScore {
		o.QualityScore = score
	}
	o.UpdatedAt = now
}
