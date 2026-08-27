package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/protocol"
)



func TestBug11_StreamHonorsCanceledContext(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background()); cancel()
    _, err := protocol.NewStream(0, time.Minute).Read(ctx, strings.NewReader(""), time.Now())
    if !errors.Is(err, context.Canceled) { t.Fatalf("err=%v", err) }
}
