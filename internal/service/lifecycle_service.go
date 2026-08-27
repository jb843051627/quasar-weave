package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (l *Lab) IssueCommand(ctx context.Context, operator, kind, target string, args []string) (model.OperatorCommand, error) {
	if operator == "" || kind == "" || target == "" {
		return model.OperatorCommand{}, fmt.Errorf("operator, kind and target are required")
	}
	command := model.OperatorCommand{ID: l.nextID("command"), Operator: operator, Kind: kind, TargetID: target, Arguments: append([]string(nil), args...), RequestedAt: l.Now()}
	if err := l.store.SaveCommand(ctx, command); err != nil {
		return model.OperatorCommand{}, err
	}
	if err := l.store.SaveAudit(ctx, model.AuditEntry{ID: l.nextID("audit"), Subject: target, Action: "command.requested", Actor: operator, After: kind, OccurredAt: l.Now()}); err != nil {
		return model.OperatorCommand{}, err
	}
	return command, nil
}

func (l *Lab) CompleteCommand(ctx context.Context, id, result, failure string) (model.OperatorCommand, error) {
	command, err := l.store.CompleteCommand(ctx, id, result, failure, l.Now())
	if err != nil {
		return model.OperatorCommand{}, err
	}
	action := "command.completed"
	if failure != "" {
		action = "command.failed"
	}
	if err := l.store.SaveAudit(ctx, model.AuditEntry{ID: l.nextID("audit"), Subject: command.TargetID, Action: action, Actor: command.Operator, Before: command.Kind, After: result, OccurredAt: l.Now()}); err != nil {
		return model.OperatorCommand{}, err
	}
	return command, nil
}

func (l *Lab) Timeline(ctx context.Context, observationID string) (model.Timeline, error) {
	return l.store.BuildTimeline(ctx, observationID)
}

func (l *Lab) RecordTransition(ctx context.Context, observationID, actor string, from, to model.ObservationStatus, reason string) error {
	history := model.StateHistory{ObservationID: observationID, Actor: actor, From: from, To: to, Reason: reason, At: l.Now()}
	if !history.DescribesTransition() {
		return fmt.Errorf("invalid state history")
	}
	return l.store.SaveAudit(ctx, model.AuditEntry{ID: l.nextID("audit"), Subject: observationID, Action: "observation.transition", Actor: actor, Before: string(from), After: string(to), OccurredAt: history.At})
}

func (l *Lab) ExpireCommands(ctx context.Context, before time.Time) (int, error) {
	commands, err := l.store.ListCommands(ctx, "")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, command := range commands {
		if !command.Complete() && command.RequestedAt.Before(before) {
			count++
		}
	}
	return count, nil
}
