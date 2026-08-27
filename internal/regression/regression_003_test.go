package service_test

import (
	"context"
	"sync"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
	"github.com/jb843051627/quasar-weave/internal/model"
)

func setupBug03(t *testing.T) (*service.Lab, *store.Store, context.Context) {
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

func TestBug03_ConcurrentFrameResultsKeepEveryCount(t *testing.T) {
    lab, _, ctx := setupBug03(t)
    const workers = 32
    var group sync.WaitGroup
    for i := 0; i < workers; i++ {
        group.Add(1)
        go func() { defer group.Done(); if _, err := lab.RecordFrameResult(ctx, "obs-001", 0.8); err != nil { t.Error(err) } }()
    }
    group.Wait()
    observation, err := lab.GetObservation(ctx, "obs-001")
    if err != nil { t.Fatal(err) }
    if observation.ReceivedFrames != workers { t.Fatalf("received=%d want=%d", observation.ReceivedFrames, workers) }
}
