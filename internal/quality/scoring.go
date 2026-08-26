package quality

import (
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type ScoreBook struct {
	results []model.QualityResult
}

func NewScoreBook() *ScoreBook {
	return &ScoreBook{results: make([]model.QualityResult, 0)}
}

func (s *ScoreBook) Add(result model.QualityResult) {
	s.results = append(s.results, result)
}

func (s *ScoreBook) Passed() int {
	count := 0
	for _, result := range s.results {
		if result.Passed {
			count++
		}
	}
	return count
}

func (s *ScoreBook) Average() float64 {
	if len(s.results) == 0 {
		return 0
	}
	var total float64
	for _, result := range s.results {
		total += result.Score
	}
	return total / float64(len(s.results))
}

func (s *ScoreBook) Latest() []model.QualityResult {
	result := append([]model.QualityResult(nil), s.results...)
	sort.Slice(result, func(i, j int) bool { return result[i].EvaluatedAt.After(result[j].EvaluatedAt) })
	return result
}
