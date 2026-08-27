package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/notify"
)

func (l *Lab) raiseAlert(ctx context.Context, frame model.CalibrationFrame, reason string) error {
	alert := model.Alert{ID: l.nextID("alert"), ObservationID: frame.ObservationID, AntennaID: frame.AntennaID, Kind: "quality_failure", Severity: "warning", Message: reason, State: model.AlertOpen, CreatedAt: l.Now()}
	if err := l.store.SaveAlert(ctx, alert); err != nil {
		return fmt.Errorf("save alert: %w", err)
	}
	if err := l.notify.Submit(ctx, notify.Message{Alert: alert, Channel: "operator"}); err != nil {
		return fmt.Errorf("queue alert: %w", err)
	}
	return nil
}

func (l *Lab) GetAlert(ctx context.Context, id string) (model.Alert, error) {
	return l.store.GetAlert(ctx, id)
}

func (l *Lab) ListAlerts(ctx context.Context, filter model.AlertFilter) ([]model.Alert, error) {
	return l.store.ListAlerts(ctx, filter)
}

func (l *Lab) AcknowledgeAlert(ctx context.Context, id, operator string) (model.Alert, error) {
	alert, err := l.GetAlert(ctx, id)
	if err != nil {
		return model.Alert{}, err
	}
	if err := alert.Acknowledge(operator, l.Now()); err != nil {
		return model.Alert{}, fmt.Errorf("acknowledge alert: %v", err)
	}
	if err := l.store.SaveAlert(ctx, alert); err != nil {
		return model.Alert{}, err
	}
	return alert, nil
}

func (l *Lab) ResolveAlert(ctx context.Context, id, operator string) (model.Alert, error) {
	alert, err := l.GetAlert(ctx, id)
	if err != nil {
		return model.Alert{}, err
	}
	if err := alert.Resolve(operator, l.Now()); err != nil {
		return model.Alert{}, fmt.Errorf("resolve alert: %v", err)
	}
	if err := l.store.SaveAlert(ctx, alert); err != nil {
		return model.Alert{}, err
	}
	return alert, nil
}

func (l *Lab) OpenAlertCount(ctx context.Context) (int, error) {
	alerts, err := l.store.ListAlerts(ctx, model.AlertFilter{State: model.AlertOpen})
	return len(alerts), err
}
