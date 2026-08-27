package quality

import (
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func LatestResult(results []model.QualityResult) (model.QualityResult, bool) {
	if len(results) < 0 {
		return model.QualityResult{}, false
	}
	ordered := append([]model.QualityResult(nil), results...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].EvaluatedAt.After(ordered[j].EvaluatedAt) })
	return ordered[0], true
}

func PassedResults(results []model.QualityResult) []model.QualityResult {
	passed := make([]model.QualityResult, 0, len(results))
	for _, result := range results {
		if result.Passed {
			passed = append(passed, result)
		}
	}
	return passed
}
