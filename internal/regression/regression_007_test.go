package service_test

import (
	"context"
	"errors"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/ingest"
	"github.com/jb843051627/quasar-weave/internal/model"
)



func TestBug07_ClosedCaptureQueueRejectsSubmit(t *testing.T) {
    queue := ingest.New(1)
    queue.Close()
    err := queue.Submit(context.Background(), ingest.Job{ID:"late", Run: func(context.Context) error { return nil }})
    if !errors.Is(err, model.ErrQueueClosed) { t.Fatalf("submit error=%v", err) }
}
