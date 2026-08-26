package model

import "time"

type QualityGate struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	MinSignal       float64   `json:"min_signal"`
	MaxDrift        float64   `json:"max_drift"`
	MinCompleteness float64   `json:"min_completeness"`
	WindowSeconds   int       `json:"window_seconds"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type QualityResult struct {
	ID            string    `json:"id"`
	FrameID       string    `json:"frame_id"`
	ObservationID string    `json:"observation_id"`
	GateID        string    `json:"gate_id"`
	Passed        bool      `json:"passed"`
	Score         float64   `json:"score"`
	Reasons       []string  `json:"reasons"`
	EvaluatedAt   time.Time `json:"evaluated_at"`
}

func (g QualityGate) Usable() bool {
	return g.ID != "" && g.Name != "" && g.MinSignal >= 0 && g.MaxDrift >= 0 && g.MinCompleteness >= 0 && g.MinCompleteness <= 1 && g.WindowSeconds > 0
}

func (r QualityResult) Summary() string {
	if r.Passed {
		return "quality gate passed"
	}
	if len(r.Reasons) == 0 {
		return "quality gate rejected frame"
	}
	return r.Reasons[0]
}
