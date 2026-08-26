package model

import (
	"math"
	"sort"
	"time"
)

type TelemetryPoint struct {
	ID         string    `json:"id"`
	AntennaID  string    `json:"antenna_id"`
	Name       string    `json:"name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	CapturedAt time.Time `json:"captured_at"`
}

func (p TelemetryPoint) Finite() bool {
	return p.ID != "" && p.AntennaID != "" && p.Name != "" && p.Unit != "" && !math.IsNaN(p.Value) && !math.IsInf(p.Value, 0) && !p.CapturedAt.IsZero()
}

type TelemetrySeries struct {
	AntennaID string           `json:"antenna_id"`
	Name      string           `json:"name"`
	Points    []TelemetryPoint `json:"points"`
}

func (s TelemetrySeries) Ordered() []TelemetryPoint {
	result := append([]TelemetryPoint(nil), s.Points...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].CapturedAt.Before(result[j].CapturedAt) })
	return result
}

func (s TelemetrySeries) Average() float64 {
	if len(s.Points) == 0 {
		return 0
	}
	var total float64
	for _, point := range s.Points {
		total += point.Value
	}
	return total / float64(len(s.Points))
}

func (s TelemetrySeries) MinMax() (float64, float64, bool) {
	if len(s.Points) == 0 {
		return 0, 0, false
	}
	min, max := s.Points[0].Value, s.Points[0].Value
	for _, point := range s.Points[1:] {
		if point.Value < min {
			min = point.Value
		}
		if point.Value > max {
			max = point.Value
		}
	}
	return min, max, true
}

type TelemetryWindow struct {
	AntennaID string    `json:"antenna_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Points    int       `json:"points"`
	Average   float64   `json:"average"`
	Minimum   float64   `json:"minimum"`
	Maximum   float64   `json:"maximum"`
}

func (w TelemetryWindow) Duration() time.Duration { return w.End.Sub(w.Start) }
