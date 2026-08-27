package model

import "time"

type RetryPlan struct {
	ID            string     `json:"id"`
	ObservationID string     `json:"observation_id"`
	Attempt       int        `json:"attempt"`
	MaxAttempts   int        `json:"max_attempts"`
	NextAt        time.Time  `json:"next_at"`
	Reason        string     `json:"reason"`
	State         RetryState `json:"state"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (r RetryPlan) Exhausted() bool {
	return r.Attempt >= r.MaxAttempts
}

func (r *RetryPlan) Advance(now time.Time, delay time.Duration) {
	r.Attempt++
	r.NextAt = now.Add(delay)
	r.UpdatedAt = now
	if r.Exhausted() {
		r.State = RetryFinished
	}
}
