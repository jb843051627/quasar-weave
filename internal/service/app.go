package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jb843051627/quasar-weave/internal/clock"
	"github.com/jb843051627/quasar-weave/internal/ingest"
	"github.com/jb843051627/quasar-weave/internal/metrics"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/notify"
	"github.com/jb843051627/quasar-weave/internal/quality"
	"github.com/jb843051627/quasar-weave/internal/store"
)

type Lab struct {
	store     *store.Store
	clock     clock.Clock
	evaluator *quality.Evaluator
	pipeline  *ingest.Pipeline
	notify    *notify.Dispatcher
	metrics   metrics.Counters
	telemetry *metrics.Registry
	sequence  atomic.Uint64
	closeOnce sync.Once
	stateMu   sync.Mutex
	frameMu   sync.Mutex
	closed    chan struct{}
}

func NewLab(repository *store.Store) *Lab {
	sink := notify.NewMemorySink()
	return NewLabWith(repository, clock.System{}, notify.NewDispatcher(sink, 32))
}

func NewLabWith(repository *store.Store, c clock.Clock, dispatcher *notify.Dispatcher) *Lab {
	if c == nil {
		c = clock.System{}
	}
	if dispatcher == nil {
		dispatcher = notify.NewDispatcher(notify.NewMemorySink(), 32)
	}
	lab := &Lab{
		store:     repository,
		clock:     c,
		evaluator: quality.NewEvaluator(c.Now),
		pipeline:  ingest.NewPipeline(64, 2),
		notify:    dispatcher,
		closed:    make(chan struct{}),
		telemetry: metrics.NewRegistry(),
	}
	return lab
}

func (l *Lab) Close() {
	l.closeOnce.Do(func() {
		close(l.closed)
		l.pipeline.Close()
		l.notify.Close()
	})
}

func (l *Lab) Store() *store.Store { return l.store }

func (l *Lab) TelemetryMetric(name string) float64 { return l.telemetry.Get(name) }

func (l *Lab) Now() time.Time { return l.clock.Now().UTC() }

func (l *Lab) nextID(prefix string) string {
	sequence := l.sequence.Add(1)
	return fmt.Sprintf("%s-%d-%d", prefix, l.Now().UnixNano(), sequence)
}

func (l *Lab) record(ctx context.Context, subject, action, payload string) error {
	return l.store.Event(ctx, subject, action, payload)
}

func (l *Lab) ensureOpen() error {
	select {
	case <-l.closed:
		return model.ErrQueueClosed
	default:
		return nil
	}
}
