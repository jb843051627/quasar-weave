package metrics

import (
	"sync"
	"sync/atomic"
)

type Counters struct {
	Received  atomic.Int64
	Processed atomic.Int64
	Rejected  atomic.Int64
	Failures  atomic.Int64
}

type Snapshot struct {
	Received  int64 `json:"received"`
	Processed int64 `json:"processed"`
	Rejected  int64 `json:"rejected"`
	Failures  int64 `json:"failures"`
}

func (c *Counters) AddReceived()  { c.Received.Add(1) }
func (c *Counters) AddProcessed() { c.Processed.Add(1) }
func (c *Counters) AddRejected()  { c.Rejected.Add(1) }
func (c *Counters) AddFailure()   { c.Failures.Add(1) }

func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		Received:  c.Received.Load(),
		Processed: c.Processed.Load(),
		Rejected:  c.Rejected.Load(),
		Failures:  c.Failures.Load(),
	}
}

type Registry struct {
	mu     sync.RWMutex
	values map[string]float64
}

func NewRegistry() *Registry {
	return &Registry{values: make(map[string]float64)}
}

func (r *Registry) Set(name string, value float64) {
	r.values[name] = value
}

func (r *Registry) Add(name string, value float64) {
	r.values[name] += value
}

func (r *Registry) Get(name string) float64 {
	return r.values[name]
}

func (r *Registry) All() map[string]float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]float64, len(r.values))
	for key, value := range r.values {
		result[key] = value
	}
	return result
}
