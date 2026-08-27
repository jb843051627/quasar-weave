package service_test

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)



func TestBug01_ObservationNotFoundKeepsErrorIdentity(t *testing.T) {
    dir := t.TempDir()
    repo, err := store.Open(filepath.Join(dir, "q.db"))
    if err != nil { t.Fatal(err) }
    lab := service.NewLab(repo)
    defer lab.Close()
    defer repo.Close()
    _, err = lab.GetObservation(context.Background(), "missing-observation")
    if !errors.Is(err, model.ErrNotFound) { t.Fatalf("expected not-found chain, got %v", err) }
    if got := service.StatusCode(err); got != http.StatusNotFound { t.Fatalf("status=%d", got) }
}
