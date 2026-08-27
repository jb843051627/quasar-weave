package quality

import (
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Aggregate struct {
	ObservationID string  `json:"observation_id"`
	Count         int     `json:"count"`
	Passed        int     `json:"passed"`
	Failed        int     `json:"failed"`
	AverageScore  float64 `json:"average_score"`
	WorstScore    float64 `json:"worst_score"`
}

func AggregateResults(results []model.QualityResult) Aggregate {
	result := Aggregate{}
	if len(results) == 0 {
		return result
	}
	result.ObservationID = results[0].ObservationID
	result.WorstScore = results[0].Score
	for _, item := range results {
		result.Count++
		result.AverageScore += item.Score
		if item.Passed {
			result.Passed++
		} else {
			result.Failed++
		}
		if item.Score < result.WorstScore {
			result.WorstScore = item.Score
		}
	}
	result.AverageScore /= float64(result.Count)
	return result
}

func SortResults(results []model.QualityResult) []model.QualityResult {
	copyOf := results
	sort.SliceStable(copyOf, func(i, j int) bool { return copyOf[i].EvaluatedAt.Before(copyOf[j].EvaluatedAt) })
	return copyOf
}
