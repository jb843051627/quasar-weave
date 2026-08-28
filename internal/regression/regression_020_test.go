package service_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/protocol"
)



type cancelReader struct { cancel context.CancelFunc; sent bool }

func (r *cancelReader) Read(p []byte) (int, error) {
    if r.sent { return 0, io.EOF }
    r.sent = true
    line := "frame|obs|ant|1|" + time.Now().UTC().Format(time.RFC3339Nano) + "|0.8|0.1|0.9|hash\n"
    copy(p, line)
    r.cancel()
    return len(line), nil
}

func TestBug20_CanceledStreamStopsBeforeNextFrame(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    _, err := protocol.NewStream(0, time.Minute).Read(ctx, &cancelReader{cancel: cancel}, time.Now().UTC())
    if !errors.Is(err, context.Canceled) { t.Fatalf("err=%v", err) }
}
