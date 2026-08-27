package service

import (
	"context"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (l *Lab) Health(ctx context.Context) (model.HealthSummary, error) {
	return l.store.BuildHealth(ctx, l.Now(), 5*time.Minute)
}

func (l *Lab) QueueDepth() int { return l.pipeline.Pending() }

func (l *Lab) Metrics() map[string]int64 {
	snapshot := l.metrics.Snapshot()
	return map[string]int64{"received": snapshot.Received, "processed": snapshot.Processed, "rejected": snapshot.Rejected, "failures": snapshot.Failures}
}

func (l *Lab) Events(ctx context.Context, subject string, limit int) ([]model.Event, error) {
	return l.store.Events(ctx, subject, limit)
}

func (l *Lab) QualityResults(ctx context.Context, observationID string) ([]model.QualityResult, error) {
	return l.store.ListQualityResults(ctx, observationID)
}
