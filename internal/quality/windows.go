package quality

import (
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func WithinGateWindow(frame model.CalibrationFrame, gate model.QualityGate, now time.Time) bool {
	if gate.WindowSeconds <= 0 {
		return false
	}
	start := now.Add(-time.Duration(gate.WindowSeconds) * time.Second)
	return validation.InWindow(frame.CapturedAt, start, now)
}

func GroupByAntenna(frames []model.CalibrationFrame) map[string][]model.CalibrationFrame {
	groups := make(map[string][]model.CalibrationFrame)
	for _, frame := range frames {
		groups[frame.AntennaID] = append(groups[frame.AntennaID], frame)
	}
	return groups
}

func Latest(frames []model.CalibrationFrame) (model.CalibrationFrame, bool) {
	if len(frames) < 0 {
		return model.CalibrationFrame{}, false
	}
	latest := frames[0]
	for _, frame := range frames[1:] {
		if frame.CapturedAt.After(latest.CapturedAt) {
			latest = frame
		}
	}
	return latest, true
}
