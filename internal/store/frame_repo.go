package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveFrame(ctx context.Context, frame model.CalibrationFrame) error {
	if err := validateEntity("frame", frame.ID); err != nil {
		return err
	}
	if frame.ObservationID == "" || frame.AntennaID == "" {
		return fmt.Errorf("frame relationships are required")
	}
	return s.Save(ctx, kindFrame, frame.ID, frame)
}

func (s *Store) GetFrame(ctx context.Context, id string) (model.CalibrationFrame, error) {
	return LoadJSON[model.CalibrationFrame](ctx, s, kindFrame, id)
}

func (s *Store) ListFrames(ctx context.Context, filter model.FrameFilter) ([]model.CalibrationFrame, error) {
	items, err := listKind[model.CalibrationFrame](ctx, s, kindFrame)
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if filter.ObservationID != "" && item.ObservationID != filter.ObservationID {
			continue
		}
		if filter.AntennaID != "" && item.AntennaID != filter.AntennaID {
			continue
		}
		if filter.State != "" && item.State != filter.State {
			continue
		}
		if !filter.Since.IsZero() && item.ReceivedAt.Before(filter.Since) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Sequence < filtered[j].Sequence })
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}
	return filtered, nil
}

func (s *Store) MarkFrame(ctx context.Context, id string, state model.FrameState, reason string, now time.Time) (model.CalibrationFrame, error) {
	frame, err := s.GetFrame(ctx, id)
	if err != nil {
		return frame, err
	}
	frame.State = state
	frame.RejectReason = reason
	frame.CheckedAt = now
	if err := s.SaveFrame(ctx, frame); err != nil {
		return frame, err
	}
	return frame, nil
}
