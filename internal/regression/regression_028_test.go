package service_test

import (
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/planner"
)



func TestBug28_FutureWindowIsNotDue(t *testing.T) {
    q:=planner.NewQueue();q.Add(planner.Window{ID:"future",Target:"target",Band:"L",Start:time.Now().Add(time.Hour),End:time.Now().Add(2*time.Hour),Priority:1});if _,ok:=q.Next(time.Now());ok{t.Fatal("future window")}
}
