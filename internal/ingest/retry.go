package ingest

import (
	"context"
	"sync"
	"time"
)

type RetryTask struct {
	ID       string
	Attempts int
	Run      func(context.Context) error
	NextAt   time.Time
}

type RetryScheduler struct {
	mu     sync.Mutex
	tasks  map[string]RetryTask
	closed bool
}

func NewRetryScheduler() *RetryScheduler {
	return &RetryScheduler{tasks: make(map[string]RetryTask)}
}

func (s *RetryScheduler) Add(task RetryTask) {
	s.mu.Lock()
	if !s.closed {
		s.tasks[task.ID] = task
	}
	s.mu.Unlock()
}

func (s *RetryScheduler) Remove(id string) {
	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()
}

func (s *RetryScheduler) Due(now time.Time) []RetryTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]RetryTask, 0)
	for _, task := range s.tasks {
		if !task.NextAt.After(now) {
			result = append(result, task)
		}
	}
	return result
}

func (s *RetryScheduler) Close() {
	s.mu.Lock()
	s.closed = true
	s.tasks = make(map[string]RetryTask)
	s.mu.Unlock()
}
