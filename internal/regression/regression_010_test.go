package service_test

import (
	"context"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)

func setupBug10(t *testing.T) (*service.Lab, *store.Store, context.Context) {
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

func TestBug10_QualityEvaluationIsIdempotent(t *testing.T) {
    lab, _, ctx := setupBug10(t)
    frame, done, err := lab.SubmitFrame(ctx, model.FrameInput{ID:"bad-frame",ObservationID:"obs-001",AntennaID:"ant-001",Sequence:1,Signal:0.01,Drift:0.8,Completeness:0.1,CapturedAt:time.Now().UTC()})
    if err != nil { t.Fatal(err) }; if done != nil { if err := <-done; err != nil { t.Fatal(err) } }
    if _, err := lab.EvaluateFrame(ctx, frame.ID); err != nil { t.Fatal(err) }
    if _, err := lab.EvaluateFrame(ctx, frame.ID); err != nil { t.Fatal(err) }
    alerts, err := lab.ListAlerts(ctx, model.AlertFilter{ObservationID:"obs-001"}); if err != nil { t.Fatal(err) }
    if len(alerts) != 1 { t.Fatalf("alerts=%d", len(alerts)) }
}
