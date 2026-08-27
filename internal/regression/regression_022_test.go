package service_test

import (
	"sync"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/metrics"
)



func TestBug22_RegistryConcurrentUpdates(t *testing.T) {
    registry := metrics.NewRegistry(); var group sync.WaitGroup; for i:=0;i<32;i++ { group.Add(1); go func(){ defer group.Done(); registry.Add("signal",1); _=registry.Get("signal") }() }; group.Wait(); if registry.Get("signal") != 32 { t.Fatalf("value=%v", registry.Get("signal")) }
}
