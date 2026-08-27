package service_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
	"github.com/jb843051627/quasar-weave/internal/store"
)



func TestBug08_MissingFrameRetainsNotFound(t *testing.T) {
    dir := t.TempDir(); repo, err := store.Open(filepath.Join(dir,"q.db")); if err != nil { t.Fatal(err) }
    lab := service.NewLab(repo); defer lab.Close(); defer repo.Close()
    _, err = lab.GetFrame(context.Background(), "missing-frame")
    if !errors.Is(err, model.ErrNotFound) { t.Fatalf("error=%v", err) }
}
