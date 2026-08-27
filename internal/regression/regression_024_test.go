package service_test

import (
	"context"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/metrics"
	"github.com/jb843051627/quasar-weave/internal/notify"
)



func TestBug24_QuerySnapshotsAreIndependent(t *testing.T) {
    sink:=notify.NewMemorySink();_=sink.Send(context.Background(),notify.Message{Channel:"operator"});m:=sink.Messages();m[0].Channel="changed";if got,_:=sink.Last();got.Channel=="changed"{t.Fatal("alias")};r:=metrics.NewRegistry();r.Set("x",1);a:=r.All();a["x"]=9;if r.Get("x")!=1{t.Fatal("map alias")}
}
