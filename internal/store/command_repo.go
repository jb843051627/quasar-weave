package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

const kindCommand = "operator_command"

func (s *Store) SaveCommand(ctx context.Context, command model.OperatorCommand) error {
	if command.ID == "" || command.Operator == "" || command.Kind == "" || command.TargetID == "" {
		return fmt.Errorf("command identity is incomplete")
	}
	return s.Save(ctx, kindCommand, command.ID, command)
}

func (s *Store) GetCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	return LoadJSON[model.OperatorCommand](ctx, s, kindCommand, id)
}

func (s *Store) ListCommands(ctx context.Context, targetID string) ([]model.OperatorCommand, error) {
	items, err := listKind[model.OperatorCommand](ctx, s, kindCommand)
	if err != nil {
		return nil, err
	}
	result := items[:0]
	for _, item := range items {
		if targetID == "" || item.TargetID == targetID {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RequestedAt.Before(result[j].RequestedAt) })
	return result, nil
}

func (s *Store) CompleteCommand(ctx context.Context, id, result string, errMessage string, at time.Time) (model.OperatorCommand, error) {
	command, err := s.GetCommand(ctx, id)
	if err != nil {
		return command, err
	}
	command.Result, command.ErrorMessage, command.CompletedAt = result, errMessage, at
	if err := s.SaveCommand(ctx, command); err != nil {
		return command, err
	}
	return command, nil
}
