package validation

import "github.com/jb843051627/quasar-weave/internal/model"

func ObservationTransition(from, to model.ObservationStatus) error {
	if !model.CanTransition(from, to) {
		return model.ErrInvalidState
	}
	return nil
}

func FrameState(value model.FrameState) bool {
	switch value {
	case model.FrameReceived, model.FrameQueued, model.FrameChecked, model.FrameRejected:
		return true
	default:
		return false
	}
}
