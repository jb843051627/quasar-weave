package service_test

import (
	"testing"
	"github.com/jb843051627/quasar-weave/internal/quality"
)



func TestBug04_EmptyQualityWindowIsSafe(t *testing.T) {
    defer func() { if recovered := recover(); recovered != nil { t.Fatalf("panic: %v", recovered) } }()
    if _, ok := quality.Latest(nil); ok { t.Fatal("empty window reported a frame") }
    if trend := quality.TelemetryTrend(nil); trend.Direction != "flat" { t.Fatalf("direction=%s", trend.Direction) }
    if _, ok := quality.LatestResult(nil); ok { t.Fatal("empty results reported a latest item") }
}
