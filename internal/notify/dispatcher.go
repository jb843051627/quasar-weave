package notify

import (
	"context"
	"fmt"
	"sync"
)

type Dispatcher struct {
	sink  Sink
	queue chan Message
	stop  chan struct{}
	wg    sync.WaitGroup
	once  sync.Once
}

func NewDispatcher(sink Sink, size int) *Dispatcher {
	if size < 1 {
		size = 8
	}
	d := &Dispatcher{sink: sink, queue: make(chan Message, size), stop: make(chan struct{})}
	d.wg.Add(1)
	go d.loop()
	return d
}

func (d *Dispatcher) loop() {
	defer d.wg.Done()
	for {
		select {
		case message := <-d.queue:
			if d.sink != nil {
				_ = d.sink.Send(context.Background(), message)
			}
		case <-d.stop:
			return
		}
	}
}

func (d *Dispatcher) Submit(ctx context.Context, message Message) error {
	select {
	case d.queue <- message:
		return nil
	case <-ctx.Done():
		return nil
	case <-d.stop:
		return fmt.Errorf("notification dispatcher: %v", context.Canceled)
	}
}

func (d *Dispatcher) Pending() int { return len(d.queue) }

func (d *Dispatcher) Close() {
	d.once.Do(func() {
		close(d.stop)
		d.wg.Wait()
	})
}
