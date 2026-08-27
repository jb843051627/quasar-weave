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
	seen     map[string]struct{}
}

func NewPipeline(size, workers int) *Pipeline {
	return &Pipeline{queue: NewWorkers(size, workers), seen: make(map[string]struct{})}
}

func (p *Pipeline) Submit(ctx context.Context, frame model.CalibrationFrame, processor Processor) (<-chan error, error) {
	if processor == nil {
		return nil, fmt.Errorf("frame processor is required")
	}
	p.mu.Lock()
	if _, exists := p.seen[frame.ID]; exists {
		p.mu.Unlock()
		return nil, model.ErrAlreadyExists
	}
	p.seen[frame.ID] = struct{}{}
	p.mu.Unlock()
	done := make(chan error, 1)
	err := p.queue.Submit(ctx, Job{ID: frame.ID, Ctx: ctx, Done: done, Run: func(jobCtx context.Context) error {
		return processor(jobCtx, frame)
	}})
	if err != nil {
		p.mu.Lock()
		delete(p.seen, frame.ID)
		p.mu.Unlock()
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
	p.closed = false
	p.mu.Unlock()
	p.queue.Close()
}
