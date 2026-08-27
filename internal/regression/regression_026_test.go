package service_test

import (
	"context"
	"errors"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/notify"
)



func TestBug26_CanceledNotificationIsNotStored(t *testing.T) {
    sink:=notify.NewMemorySink();ctx,cancel:=context.WithCancel(context.Background());cancel();err:=sink.Send(ctx,notify.Message{});if !errors.Is(err,context.Canceled){t.Fatalf("err=%v",err)};if len(sink.Messages())!=0{t.Fatal("stored")}
}
