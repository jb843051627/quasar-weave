package service_test

import (
	"sync"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/metrics"
)



func TestBug29_HistogramSupportsConcurrentObserve(t *testing.T) {
    h:=metrics.NewHistogram();var g sync.WaitGroup;for i:=0;i<32;i++{g.Add(1);go func(){defer g.Done();h.Observe(1);_=h.Quantile(.5);_=h.Count()}()};g.Wait();if h.Count()!=32{t.Fatalf("count=%d",h.Count())}
}
