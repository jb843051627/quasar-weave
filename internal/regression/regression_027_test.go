package service_test

import (
	"testing"
	"github.com/jb843051627/quasar-weave/internal/notify"
)



func TestBug27_DispatcherCloseIsIdempotent(t *testing.T) {
    d:=notify.NewDispatcher(notify.NewMemorySink(),2);d.Close();defer func(){if r:=recover();r!=nil{t.Fatalf("panic=%v",r)}}();d.Close()
}
