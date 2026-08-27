package model

import "time"

type Alert struct {
	ID            string     `json:"id"`
	ObservationID string     `json:"observation_id"`
	AntennaID     string     `json:"antenna_id"`
	Kind          string     `json:"kind"`
	Severity      string     `json:"severity"`
	Message       string     `json:"message"`
	State         AlertState `json:"state"`
	CreatedAt     time.Time  `json:"created_at"`
	AckAt         time.Time  `json:"ack_at"`
	ResolvedAt    time.Time  `json:"resolved_at"`
	Operator      string     `json:"operator"`
}

func (a Alert) Open() bool {
	return a.State == AlertOpen || a.State == AlertAcked
}

func (a *Alert) Acknowledge(operator string, now time.Time) error {
	if a.State != AlertOpen {
		return ErrInvalidState
	}
	a.State = AlertAcked
	a.Operator = operator
	a.AckAt = now
	return nil
}

func (a *Alert) Resolve(operator string, now time.Time) error {
	if a.State != AlertOpen && a.State != AlertAcked {
		return ErrInvalidState
	}
	a.State = AlertResolved
	a.Operator = operator
	a.ResolvedAt = now
	return nil
}
