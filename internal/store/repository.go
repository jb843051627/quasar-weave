package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
)

const (
	kindAntenna     = "antenna"
	kindObservation = "observation"
	kindFrame       = "calibration_frame"
	kindGate        = "quality_gate"
	kindResult      = "quality_result"
	kindRetry       = "retry_plan"
	kindAlert       = "alert"
	kindNote        = "operator_note"
)

func decodeOne[T any](raw []byte, kind string) (T, error) {
	var value T
	if err := json.Unmarshal(raw, &value); err != nil {
		return value, fmt.Errorf("decode %s: %w", kind, err)
	}
	return value, nil
}

func listKind[T any](ctx context.Context, s *Store, kind string) ([]T, error) {
	values := make([]T, 0)
	err := s.List(ctx, kind, func(raw []byte) error {
		value, err := decodeOne[T](raw, kind)
		if err != nil {
			return err
		}
		values = append(values, value)
		return nil
	})
	return values, err
}

func validateEntity(kind, id string) error {
	return model.ValidateID(kind, id)
}
