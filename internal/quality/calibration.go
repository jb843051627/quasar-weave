package quality

import (
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type CalibrationSummary struct {
	Frames       int     `json:"frames"`
	Passed       int     `json:"passed"`
	Rejected     int     `json:"rejected"`
	MeanSignal   float64 `json:"mean_signal"`
	MeanDrift    float64 `json:"mean_drift"`
	MeanCoverage float64 `json:"mean_coverage"`
	Score        float64 `json:"score"`
}

func SummarizeFrames(frames []model.CalibrationFrame) CalibrationSummary {
	result := CalibrationSummary{Frames: len(frames)}
	if len(frames) == 0 {
		return result
	}
	for _, frame := range frames {
		result.MeanSignal += frame.Signal
		result.MeanDrift += frame.Drift
		result.MeanCoverage += frame.Completeness
		if frame.State == model.FrameChecked {
			result.Passed++
		}
		if frame.State == model.FrameRejected {
			result.Rejected++
		}
	}
	count := float64(len(frames))
	result.MeanSignal /= count
	result.MeanDrift /= count
	result.MeanCoverage /= count
	result.Score = result.MeanSignal*0.4 + (1-result.MeanDrift)*0.2 + result.MeanCoverage*0.4
	return result
}

func SelectRepresentative(frames []model.CalibrationFrame, count int) ([]model.CalibrationFrame, error) {
	if count <= 0 {
		return nil, fmt.Errorf("representative count must be positive")
	}
	ordered := append([]model.CalibrationFrame(nil), frames...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Sequence != ordered[j].Sequence {
			return ordered[i].Sequence < ordered[j].Sequence
		}
		return ordered[i].CapturedAt.Before(ordered[j].CapturedAt)
	})
	if len(ordered) <= count {
		return ordered, nil
	}
	result := make([]model.CalibrationFrame, 0, count)
	step := float64(len(ordered)-1) / float64(count-1)
	for i := 0; i < count; i++ {
		result = append(result, ordered[int(float64(i)*step)])
	}
	return result, nil
}

func CompareSummaries(before, after CalibrationSummary) float64 {
	return (after.Score - before.Score) + float64(after.Passed-before.Passed)*0.01 - float64(after.Rejected-before.Rejected)*0.01
}
