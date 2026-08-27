package metrics

import "sync"

type Rolling struct {
	mu     sync.Mutex
	values []float64
	limit  int
}

func NewRolling(limit int) *Rolling {
	if limit < 1 {
		limit = 1
	}
	return &Rolling{limit: limit, values: make([]float64, 0, limit)}
}

func (r *Rolling) Add(value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values = append(r.values, value)
	if len(r.values) > r.limit {
		r.values = append([]float64(nil), r.values[len(r.values)-r.limit:]...)
	}
}

func (r *Rolling) Values() []float64 {
	return r.values
}

func (r *Rolling) Average() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.values) == 0 {
		return 0
	}
	var total float64
	for _, value := range r.values {
		total += value
	}
	return total / float64(len(r.values))
}
