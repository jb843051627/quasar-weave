package model

import (
	"fmt"
	"strings"
	"time"
)

type ObservationRule struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	AllowedBands    []string      `json:"allowed_bands"`
	MinimumFrames   int           `json:"minimum_frames"`
	MaximumDuration time.Duration `json:"maximum_duration"`
	RequiresNote    bool          `json:"requires_note"`
}

func (r ObservationRule) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("rule identity is required")
	}
	if r.MinimumFrames < 1 || r.MaximumDuration <= 0 {
		return fmt.Errorf("rule limits are invalid")
	}
	if len(r.AllowedBands) == 0 {
		return fmt.Errorf("at least one band is required")
	}
	return nil
}

func (r ObservationRule) AllowsBand(band string) bool {
	for _, allowed := range r.AllowedBands {
		if strings.EqualFold(allowed, band) {
			return true
		}
	}
	return false
}

func (r ObservationRule) Satisfied(observation Observation, elapsed time.Duration) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if observation.ReceivedFrames < r.MinimumFrames {
		return fmt.Errorf("observation has %d frames, need %d", observation.ReceivedFrames, r.MinimumFrames)
	}
	if elapsed > r.MaximumDuration {
		return fmt.Errorf("observation exceeded maximum duration")
	}
	if observation.Status != ObservationQualified && observation.Status != ObservationArchived {
		return fmt.Errorf("observation is not qualified")
	}
	return nil
}

type StateHistory struct {
	ObservationID string            `json:"observation_id"`
	From          ObservationStatus `json:"from"`
	To            ObservationStatus `json:"to"`
	At            time.Time         `json:"at"`
	Actor         string            `json:"actor"`
	Reason        string            `json:"reason"`
}

func (h StateHistory) DescribesTransition() bool {
	return h.ObservationID != "" && h.From != h.To && CanTransition(h.From, h.To) && !h.At.IsZero()
}
