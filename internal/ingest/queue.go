package ingest

import (
	"context"
	"sync"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Job struct {
	ID   string
	Ctx  context.Context
	Run  func(context.Context) error
	Done chan error
}

type Queue struct {
	jobs   chan Job
	stop   chan struct{}
	done   chan struct{}
	once   sync.Once
	wg     sync.WaitGroup
	mu     sync.RWMutex
	closed bool
}

func New(size int) *Queue { return NewWorkers(size, 1) }

func NewWorkers(size, workers int) *Queue {
	if size < 1 {
		size = 16
	}
	if workers < 1 {
		workers = 1
	}
	q := &Queue{jobs: make(chan Job, size), stop: make(chan struct{}), done: make(chan struct{})}
	q.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go q.loop()
	}
	return q
}

func (q *Queue) loop() {
	defer q.wg.Done()
	for {
		select {
		case job := <-q.jobs:
			if job.Run == nil {
				q.finish(job, model.ErrQueueClosed)
				continue
			}
			jobCtx := job.Ctx
			if jobCtx == nil {
				jobCtx = context.Background()
			}
			err := job.Run(jobCtx)
			q.finish(job, err)
		case <-q.stop:
			return
		}
	}
}

func (q *Queue) finish(job Job, err error) {
	if job.Done == nil {
		return
	}
	select {
	case job.Done <- err:
	default:
	}
}

func (q *Queue) Submit(ctx context.Context, job Job) error {
	q.mu.RLock()
	closed := q.closed
	q.mu.RUnlock()
	if closed {
		return model.ErrQueueClosed
	}
	if job.Ctx == nil {
		job.Ctx = ctx
	}
	select {
	case q.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-q.stop:
		return model.ErrQueueClosed
	}
}

func (q *Queue) Len() int { return len(q.jobs) }

func (q *Queue) Close() {
	q.once.Do(func() {
		q.closed = true
		q.wg.Wait()
		close(q.stop)
		close(q.done)
	})
}

func (q *Queue) Done() <-chan struct{} { return q.done }
