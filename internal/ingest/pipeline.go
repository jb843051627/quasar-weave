package ingest

import (
	"context"
	"fmt"
	"sync"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Processor func(context.Context, model.CalibrationFrame) error

type Pipeline struct {
	queue    *Queue
	mu       sync.RWMutex
	closed   bool
	accepted int
	failed   int
}

func NewPipeline(size, workers int) *Pipeline {
	return &Pipeline{queue: NewWorkers(size, workers)}
}

func (p *Pipeline) Submit(ctx context.Context, frame model.CalibrationFrame, processor Processor) (<-chan error, error) {
	if processor == nil {
		return nil, fmt.Errorf("frame processor is required")
	}
	done := make(chan error, 1)
	err := p.queue.Submit(ctx, Job{ID: frame.ID, Ctx: ctx, Done: done, Run: func(jobCtx context.Context) error {
		return processor(jobCtx, frame)
	}})
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.accepted++
	p.mu.Unlock()
	return done, nil
}

func (p *Pipeline) Pending() int { return p.queue.Len() }

func (p *Pipeline) Accepted() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.accepted
}

func (p *Pipeline) Close() {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()
	p.queue.Close()
}
