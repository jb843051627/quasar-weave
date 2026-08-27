package service_test

import (
	"testing"
	"github.com/jb843051627/quasar-weave/internal/quality"
)



func TestBug23_EmptyQualityInputsAreSafe(t *testing.T) {
    defer func(){ if recovered:=recover(); recovered!=nil { t.Fatalf("panic=%v",recovered) } }(); if _,ok:=quality.Latest(nil);ok{t.Fatal("frame")};if _,ok:=quality.LatestResult(nil);ok{t.Fatal("result")};if quality.TelemetryTrend(nil).Direction!="flat"{t.Fatal("trend")}
}
