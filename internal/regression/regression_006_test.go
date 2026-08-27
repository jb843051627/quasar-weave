package service_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/store"
)



func TestBug06_TelemetrySummaryHonorsWindow(t *testing.T) {
    dir := t.TempDir()
    repo, err := store.Open(filepath.Join(dir, "q.db")); if err != nil { t.Fatal(err) }
    defer repo.Close()
    now := time.Now().UTC()
    for _, point := range []model.TelemetryPoint{{ID:"inside",AntennaID:"ant",Name:"power",Value:2,Unit:"u",CapturedAt:now.Add(-time.Minute)},{ID:"late",AntennaID:"ant",Name:"power",Value:9,Unit:"u",CapturedAt:now.Add(time.Minute)}} { if err := repo.SaveTelemetry(context.Background(), point); err != nil { t.Fatal(err) } }
    points, err := repo.ListTelemetry(context.Background(), "ant", "power", now.Add(-2*time.Minute), now)
    if err != nil { t.Fatal(err) }
    if len(points) != 1 || points[0].ID != "inside" { t.Fatalf("points=%v", points) }
}
