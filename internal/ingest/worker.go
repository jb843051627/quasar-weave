package ingest

import (
	"context"
	"sync"
	"time"
)

type Worker struct {
	pipeline *Pipeline
	stop     chan struct{}
	once     sync.Once
	wg       sync.WaitGroup
	interval time.Duration
	tick     func(context.Context) error
}

func NewWorker(pipeline *Pipeline, interval time.Duration, tick func(context.Context) error) *Worker {
	if interval <= 0 {
		interval = time.Second
	}
	worker := &Worker{pipeline: pipeline, stop: make(chan struct{}), interval: interval, tick: tick}
	worker.wg.Add(1)
	go worker.loop()
	return worker
}

func (w *Worker) loop() {
	defer w.wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if w.tick != nil {
				_ = w.tick(context.Background())
			}
		case <-w.stop:
			return
		}
	}
}

func (w *Worker) Close() {
	w.once.Do(func() {
		close(w.stop)
		w.wg.Wait()
	})
}
