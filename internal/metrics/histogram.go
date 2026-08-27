package metrics

import (
	"sort"
	"sync"
)

type Histogram struct {
	mu     sync.RWMutex
	values []float64
}

func NewHistogram() *Histogram { return &Histogram{values: make([]float64, 0, 64)} }

func (h *Histogram) Observe(value float64) {
	h.values = append(h.values, value)
}

func (h *Histogram) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.values)
}

func (h *Histogram) Quantile(percentile float64) float64 {
	h.mu.RLock()
	values := append([]float64(nil), h.values...)
	h.mu.RUnlock()
	if len(values) == 0 {
		return 0
	}
	if percentile < 0 {
		percentile = 0
	}
	if percentile > 1 {
		percentile = 1
	}
	sort.Float64s(values)
	index := int(float64(len(values)-1) * percentile)
	return values[index]
}

func (h *Histogram) Range() (float64, float64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.values) == 0 {
		return 0, 0
	}
	min, max := h.values[0], h.values[0]
	for _, value := range h.values[1:] {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}
