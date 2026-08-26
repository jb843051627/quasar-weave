package planner

import (
	"container/heap"
	"sync"
	"time"
)

type item struct {
	Window Window
	index  int
}

type priorityQueue []*item

func (p priorityQueue) Len() int { return len(p) }
func (p priorityQueue) Less(i, j int) bool {
	if p[i].Window.Priority != p[j].Window.Priority {
		return p[i].Window.Priority > p[j].Window.Priority
	}
	return p[i].Window.Start.Before(p[j].Window.Start)
}
func (p priorityQueue) Swap(i, j int) { p[i], p[j] = p[j], p[i]; p[i].index = i; p[j].index = j }
func (p *priorityQueue) Push(value any) {
	entry := value.(*item)
	entry.index = len(*p)
	*p = append(*p, entry)
}
func (p *priorityQueue) Pop() any {
	old := *p
	n := len(old)
	entry := old[n-1]
	old[n-1] = nil
	entry.index = -1
	*p = old[:n-1]
	return entry
}

type Queue struct {
	mu    sync.Mutex
	items priorityQueue
}

func NewQueue() *Queue {
	q := &Queue{items: make(priorityQueue, 0)}
	heap.Init(&q.items)
	return q
}

func (q *Queue) Add(window Window) {
	heap.Push(&q.items, &item{Window: window})
}

func (q *Queue) Next(now time.Time) (Window, bool) {
	if len(q.items) == 0 {
		return Window{}, false
	}
	entry := q.items[0]
	if entry.Window.Start.After(now) {
		return Window{}, false
	}
	return heap.Pop(&q.items).(*item).Window, true
}

func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }

func (q *Queue) Snapshot() []Window {
	q.mu.Lock()
	defer q.mu.Unlock()
	result := make([]Window, 0, len(q.items))
	for _, entry := range q.items {
		result = append(result, entry.Window)
	}
	return SortWindows(result)
}
