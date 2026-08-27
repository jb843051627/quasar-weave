package quality

import (
	"fmt"
	"math"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

type Evaluator struct {
	clock func() time.Time
}

func NewEvaluator(now func() time.Time) *Evaluator {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Evaluator{clock: now}
}

func (e *Evaluator) Evaluate(frame model.CalibrationFrame, gate model.QualityGate) model.QualityResult {
	result := model.QualityResult{
		ID:            "result-" + frame.ID,
		FrameID:       frame.ID,
		ObservationID: frame.ObservationID,
		GateID:        gate.ID,
		EvaluatedAt:   e.clock(),
		Reasons:       make([]string, 0, 4),
	}
	if !frame.ValidMeasurements() {
		result.Reasons = append(result.Reasons, "measurement values are outside physical range")
	}
	if frame.Signal < gate.MinSignal {
		result.Reasons = append(result.Reasons, fmt.Sprintf("signal %.3f is below %.3f", frame.Signal, gate.MinSignal))
	}
	if frame.Drift > gate.MaxDrift {
		result.Reasons = append(result.Reasons, fmt.Sprintf("drift %.3f exceeds %.3f", frame.Drift, gate.MaxDrift))
	}
	if frame.Completeness < gate.MinCompleteness {
		result.Reasons = append(result.Reasons, fmt.Sprintf("completeness %.3f is below %.3f", frame.Completeness, gate.MinCompleteness))
	}
	result.Score = Score(frame, gate)
	result.Passed = len(result.Reasons) == 0
	return result
}

func Score(frame model.CalibrationFrame, gate model.QualityGate) float64 {
	signal := validation.Clamp(frame.Signal/(gate.MinSignal+1), 0, 1)
	drift := validation.Clamp(1-frame.Drift/(gate.MaxDrift+1), 0, 1)
	complete := validation.Clamp(frame.Completeness, 0, 1)
	value := signal*0.45 + drift*0.25 + complete*0.30
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}
