package quality

import (
	"math"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Anomaly struct {
	SampleID string  `json:"sample_id"`
	Kind     string  `json:"kind"`
	Value    float64 `json:"value"`
	Baseline float64 `json:"baseline"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason"`
}

func DetectAnomalies(samples []model.TelemetryPoint, threshold float64) []Anomaly {
	if len(samples) == 0 {
		return nil
	}
	baseline := meanTelemetry(samples)
	deviation := standardDeviation(samples, baseline)
	if deviation == 0 {
		deviation = 1
	}
	result := make([]Anomaly, 0)
	for _, sample := range samples {
		score := math.Abs(sample.Value-baseline) / deviation
		if score >= threshold {
			result = append(result, Anomaly{SampleID: sample.ID, Kind: sample.Name, Value: sample.Value, Baseline: baseline, Score: score, Reason: "sample deviates from window baseline"})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Score > result[j].Score })
	return result
}

func meanTelemetry(samples []model.TelemetryPoint) float64 {
	if len(samples) == 0 {
		return 0
	}
	var total float64
	for _, sample := range samples {
		total += sample.Value
	}
	return total / float64(len(samples))
}

func standardDeviation(samples []model.TelemetryPoint, mean float64) float64 {
	if len(samples) < 2 {
		return 0
	}
	var total float64
	for _, sample := range samples {
		delta := sample.Value - mean
		total += delta * delta
	}
	return math.Sqrt(total / float64(len(samples)-1))
}

func QualityLabel(score float64) string {
	switch {
	case score >= 5:
		return "critical"
	case score >= 3:
		return "warning"
	case score >= 2:
		return "review"
	default:
		return "normal"
	}
}
