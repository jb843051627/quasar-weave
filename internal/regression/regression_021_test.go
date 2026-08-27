package service_test

import (
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
)



func TestBug21_TerminalRetryStaysFinished(t *testing.T) {
     plan := model.RetryPlan{State:model.RetryRunning,Attempt:1,MaxAttempts:2}; plan.Advance(time.Now(), time.Second); if plan.State != model.RetryFinished { t.Fatalf("state=%s", plan.State) }
}
