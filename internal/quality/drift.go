package quality

import "math"

func DriftDelta(previous, current float64) float64 {
	return math.Abs(current - previous)
}

func Stable(values []float64, tolerance float64) bool {
	if len(values) < 2 {
		return true
	}
	for i := 1; i < len(values); i++ {
		if DriftDelta(values[i-1], values[i]) > tolerance {
			return false
		}
	}
	return true
}

func WeightedAverage(values, weights []float64) float64 {
	if len(values) == 0 || len(values) != len(weights) {
		return 0
	}
	var total, weight float64
	for i, value := range values {
		total += value * weights[i]
		weight += weights[i]
	}
	if weight == 0 {
		return 0
	}
	return total / weight
}
