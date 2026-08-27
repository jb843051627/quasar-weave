package service_test

import (
	"errors"
	"testing"
	"time"
	"github.com/jb843051627/quasar-weave/internal/model"
)



func TestBug17_AlertCannotBeAcknowledgedTwice(t *testing.T) {
    alert := model.Alert{State:model.AlertOpen}; now := time.Now(); if err := alert.Acknowledge("operator", now); err != nil { t.Fatal(err) }; if err := alert.Acknowledge("operator", now); !errors.Is(err, model.ErrInvalidState) { t.Fatalf("err=%v", err) }
}
