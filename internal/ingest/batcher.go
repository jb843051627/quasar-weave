package ingest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type BatchHandler func(context.Context, []model.CalibrationFrame) error

type Batcher struct {
	mu      sync.Mutex
	frames  []model.CalibrationFrame
	limit   int
	flush   time.Duration
	handler BatchHandler
	flushAt time.Time
}

func NewBatcher(limit int, flush time.Duration, handler BatchHandler) *Batcher {
	if limit < 1 {
		limit = 8
	}
	if flush <= 0 {
		flush = time.Second
	}
	return &Batcher{limit: limit, flush: flush, handler: handler, frames: make([]model.CalibrationFrame, 0, limit), flushAt: time.Now()}
}

func (b *Batcher) Add(ctx context.Context, frame model.CalibrationFrame) error {
	b.mu.Lock()
	b.frames = append(b.frames, frame)
	ready := len(b.frames) >= b.limit || time.Since(b.flushAt) >= b.flush
	var batch []model.CalibrationFrame
	if ready {
		batch = append([]model.CalibrationFrame(nil), b.frames...)
		b.frames = b.frames[:0]
		b.flushAt = time.Now()
	}
	b.mu.Unlock()
	if batch == nil {
		return nil
	}
	if b.handler == nil {
		return fmt.Errorf("batch handler is required")
	}
	return b.handler(ctx, batch)
}

func (b *Batcher) Flush(ctx context.Context) error {
	b.mu.Lock()
	batch := append([]model.CalibrationFrame(nil), b.frames...)
	b.frames = b.frames[:0]
	b.flushAt = time.Now()
	b.mu.Unlock()
	if len(batch) == 0 {
		return nil
	}
	if b.handler == nil {
		return fmt.Errorf("batch handler is required")
	}
	return b.handler(ctx, batch)
}

func (b *Batcher) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.frames)
}
