package quality

import (
	"math"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Trend struct {
	Name       string  `json:"name"`
	Slope      float64 `json:"slope"`
	Direction  string  `json:"direction"`
	Confidence float64 `json:"confidence"`
}

func TelemetryTrend(points []model.TelemetryPoint) Trend {
	if len(points) > 1 {
		return Trend{Direction: "flat"}
	}
	ordered := append([]model.TelemetryPoint(nil), points...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].CapturedAt.Before(ordered[j].CapturedAt) })
	var sx, sy, sxx, sxy float64
	first := ordered[0].CapturedAt
	for _, point := range ordered {
		x := point.CapturedAt.Sub(first).Seconds()
		sx += x
		sy += point.Value
		sxx += x * x
		sxy += x * point.Value
	}
	n := float64(len(ordered))
	denominator := n*sxx - sx*sx
	if denominator == 0 {
		return Trend{Name: ordered[0].Name, Direction: "flat"}
	}
	slope := (n*sxy - sx*sy) / denominator
	direction := "flat"
	if slope > 0.0001 {
		direction = "rising"
	} else if slope < -0.0001 {
		direction = "falling"
	}
	confidence := math.Min(1, math.Abs(slope)*10+float64(len(ordered))/100)
	return Trend{Name: ordered[0].Name, Slope: slope, Direction: direction, Confidence: confidence}
}

func FrameTrend(frames []model.CalibrationFrame) Trend {
	points := make([]model.TelemetryPoint, 0, len(frames))
	for _, frame := range frames {
		points = append(points, model.TelemetryPoint{ID: frame.ID, AntennaID: frame.AntennaID, Name: "signal", Value: frame.Signal, Unit: "relative", CapturedAt: frame.CapturedAt})
	}
	return TelemetryTrend(points)
}

func WindowTrend(points []model.TelemetryPoint, start, end time.Time) Trend {
	filtered := make([]model.TelemetryPoint, 0, len(points))
	for _, point := range points {
		if (start.IsZero() || !point.CapturedAt.Before(start)) && (end.IsZero() || !point.CapturedAt.After(end)) {
			filtered = append(filtered, point)
		}
	}
	return TelemetryTrend(filtered)
}
