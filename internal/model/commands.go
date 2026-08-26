package model

import (
	"strings"
	"time"
)

type OperatorCommand struct {
	ID           string    `json:"id"`
	Operator     string    `json:"operator"`
	Kind         string    `json:"kind"`
	TargetID     string    `json:"target_id"`
	Arguments    []string  `json:"arguments"`
	RequestedAt  time.Time `json:"requested_at"`
	CompletedAt  time.Time `json:"completed_at"`
	Result       string    `json:"result"`
	ErrorMessage string    `json:"error_message"`
}

func (c OperatorCommand) Complete() bool { return !c.CompletedAt.IsZero() }

func (c OperatorCommand) Successful() bool { return c.Complete() && c.ErrorMessage == "" }

func (c OperatorCommand) Summary() string {
	if c.ErrorMessage != "" {
		return "failed: " + c.ErrorMessage
	}
	if !c.Complete() {
		return "pending"
	}
	if strings.TrimSpace(c.Result) == "" {
		return "completed"
	}
	return c.Result
}

type CommandBatch struct {
	ID        string            `json:"id"`
	Commands  []OperatorCommand `json:"commands"`
	CreatedAt time.Time         `json:"created_at"`
}

func (b CommandBatch) CompletedCount() int {
	count := 0
	for _, command := range b.Commands {
		if command.Complete() {
			count++
		}
	}
	return count
}

func (b CommandBatch) FailedCount() int {
	count := 0
	for _, command := range b.Commands {
		if command.ErrorMessage != "" {
			count++
		}
	}
	return count
}

func (b CommandBatch) Done() bool {
	return len(b.Commands) > 0 && b.CompletedCount() == len(b.Commands)
}
