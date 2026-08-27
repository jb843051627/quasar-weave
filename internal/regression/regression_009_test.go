package service_test

import (
	"context"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)

func setupBug09(t *testing.T) (*service.Lab, *store.Store, context.Context) {
	t.Helper()
	dir := t.TempDir()
	repo, err := store.Open(dir + "/quasar.db")
	if err != nil { t.Fatal(err) }
	lab := service.NewLab(repo)
	ctx := context.Background()
	if _, err := lab.EnsureDefaultGate(ctx); err != nil { t.Fatal(err) }

	if _, err := lab.RegisterAntenna(ctx, model.AntennaInput{ID: "ant-001", Name: "North Dish", Station: "N1", Band: "L", Enabled: true}); err != nil { t.Fatal(err) }

	if _, err := lab.CreateObservation(ctx, model.ObservationInput{ID: "obs-001", Target: "J1939-6342", RequestedBy: "operator", ExpectedFrames: 64, GateID: "default-gate"}); err != nil { t.Fatal(err) }
	t.Cleanup(func() { lab.Close(); repo.Close() })
	return lab, repo, ctx
}

func TestBug09_EmptyObservationSummaryDoesNotPanic(t *testing.T) {
    lab, _, ctx := setupBug09(t)
    defer func() { if recovered := recover(); recovered != nil { t.Fatalf("panic: %v", recovered) } }()
    summary, err := lab.ObservationSummary(ctx, "obs-001")
    if err != nil { t.Fatal(err) }
     if summary.Frames != 0 || summary.Aggregate.Count != 0 { t.Fatalf("summary=%+v", summary) }
 }
