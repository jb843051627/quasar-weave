package planner

import (
	"math"
	"time"
)

type Estimate struct {
	Frames     int           `json:"frames"`
	Duration   time.Duration `json:"duration"`
	DataBytes  int64         `json:"data_bytes"`
	FinishAt   time.Time     `json:"finish_at"`
	Confidence float64       `json:"confidence"`
}

func EstimateRun(start time.Time, strategy Strategy, bytesPerFrame int64) Estimate {
	duration := strategy.CaptureDuration()
	confidence := 1.0
	if strategy.RetryLimit > 0 {
		confidence -= float64(strategy.RetryLimit) * 0.03
	}
	confidence = math.Max(0.1, math.Min(1, confidence))
	return Estimate{Frames: strategy.ExpectedFrames, Duration: duration, DataBytes: int64(strategy.ExpectedFrames) * bytesPerFrame, FinishAt: start.Add(duration), Confidence: confidence}
}

func MergeEstimates(items []Estimate) Estimate {
	result := Estimate{Confidence: 1}
	for _, item := range items {
		result.Frames += item.Frames
		result.Duration += item.Duration
		result.DataBytes += item.DataBytes
		if item.FinishAt.After(result.FinishAt) {
			result.FinishAt = item.FinishAt
		}
		if item.Confidence < result.Confidence {
			result.Confidence = item.Confidence
		}
	}
	return result
}
