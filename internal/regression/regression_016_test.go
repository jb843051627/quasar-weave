package service_test

import (
	"testing"
	"github.com/jb843051627/quasar-weave/internal/notify"
)



func TestBug16_EmptyNotificationHistoryIsSafe(t *testing.T) {
    sink := notify.NewMemorySink(); defer func() { if recovered := recover(); recovered != nil { t.Fatalf("panic=%v", recovered) } }(); if _, ok := sink.Last(); ok { t.Fatal("message exists") }
}
