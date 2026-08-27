package service_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/planner"
)



func TestBug19_PlannerQueueIsConcurrencySafe(t *testing.T) {
    queue := planner.NewQueue()
    now := time.Now()
    var group sync.WaitGroup
    for i := 0; i < 32; i++ {
        group.Add(1)
        go func(seq int) {
            defer group.Done()
            w := planner.Window{ID: fmt.Sprintf("w-%d", seq), Target: "target", Band: "L", Start: now.Add(time.Duration(seq) * time.Minute), End: now.Add(time.Duration(seq+1) * time.Hour), Priority: seq % 100}
            queue.Add(w)
            queue.Next(now)
        }(i)
    }
    group.Wait()
    if queue.Len() > 32 { t.Fatalf("queue corruption: len=%d", queue.Len()) }
}

func TestBug19_DisabledAntennaIsNotAllocated(t *testing.T) {
    w := planner.Window{ID:"w",Target:"target",Band:"L",Start:time.Now(),End:time.Now().Add(time.Hour),Priority:1}
    antennas := []model.Antenna{{ID:"disabled",Band:"L",Enabled:false,Status:model.AntennaReady}}
    assignments, err := planner.Allocate([]planner.Window{w}, antennas)
    if err == nil || len(assignments) != 0 { t.Fatalf("assignments=%+v err=%v", assignments, err) }
}
