package service_test

import (
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
)



func TestBug30_TimelineLastActionIsNewest(t *testing.T) {
    now:=time.Now();timeline:=model.Timeline{Events:[]model.AuditEntry{{Action:"created",OccurredAt:now.Add(-time.Hour)},{Action:"calibrated",OccurredAt:now}}};if timeline.LastAction()!="calibrated"{t.Fatalf("last=%s",timeline.LastAction())}
}
