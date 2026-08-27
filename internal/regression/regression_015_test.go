package service_test

import (
	"context"
	"errors"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)

func setupBug15(t *testing.T) (*service.Lab, *store.Store, context.Context) {
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

func TestBug15_FailedObservationCannotBeginCalibration(t *testing.T) {
    lab, _, ctx := setupBug15(t); if _, err := lab.StartObservation(ctx, "obs-001"); err != nil { t.Fatal(err) }; if _, err := lab.FailObservation(ctx, "obs-001", "bad"); err != nil { t.Fatal(err) }
    if _, err := lab.BeginCalibration(ctx, "obs-001"); !errors.Is(err, model.ErrInvalidState) { t.Fatalf("err=%v", err) }
}
