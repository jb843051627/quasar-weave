package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/planner"
)

type ObservationPlan struct {
	Window     planner.Window     `json:"window"`
	Strategy   planner.Strategy   `json:"strategy"`
	Allocation planner.Assignment `json:"allocation"`
	Estimate   planner.Estimate   `json:"estimate"`
}

func (l *Lab) PlanObservation(ctx context.Context, observationID string, start time.Time, duration time.Duration, priority int) (ObservationPlan, error) {
	observation, err := l.GetObservation(ctx, observationID)
	if err != nil {
		return ObservationPlan{}, err
	}
	if duration <= 0 {
		return ObservationPlan{}, fmt.Errorf("duration must be positive")
	}
	window := planner.Window{ID: observation.ID, Target: observation.Target, Start: start, End: start.Add(duration), Priority: priority, Band: "L"}
	if err := window.Valid(); err != nil {
		return ObservationPlan{}, err
	}
	strategy := planner.Strategy{Name: "standard-continuum", Band: window.Band, FrameInterval: time.Minute, ExpectedFrames: observation.ExpectedFrames, RetryLimit: 3, RequireHeartbeat: true}
	if err := strategy.Valid(); err != nil {
		return ObservationPlan{}, err
	}
	antennas, err := l.ListAntennas(ctx)
	if err != nil {
		return ObservationPlan{}, err
	}
	assignments, err := planner.Allocate([]planner.Window{window}, antennas)
	if err != nil {
		return ObservationPlan{}, err
	}
	if len(assignments) == 0 {
		return ObservationPlan{}, fmt.Errorf("no antenna assignment")
	}
	estimate := planner.EstimateRun(start, strategy, 4096)
	return ObservationPlan{Window: window, Strategy: strategy, Allocation: assignments[0], Estimate: estimate}, nil
}

func (l *Lab) ValidatePlan(plan ObservationPlan) error {
	if err := plan.Window.Valid(); err != nil {
		return err
	}
	if err := plan.Strategy.Valid(); err != nil {
		return err
	}
	if plan.Allocation.WindowID != plan.Window.ID {
		return fmt.Errorf("allocation window mismatch")
	}
	if len(plan.Allocation.Antenna.ID) == 0 {
		return fmt.Errorf("plan has no antenna")
	}
	return nil
}

func planInput(observation model.Observation, now time.Time) planner.Window {
	return planner.Window{ID: observation.ID, Target: observation.Target, Start: now, End: now.Add(time.Hour), Priority: 1, Band: "L"}
}
