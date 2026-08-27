package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

func (l *Lab) AddNote(ctx context.Context, input model.NoteInput) (model.OperatorNote, error) {
	if err := validation.Text(input.Operator, "operator", 2, 80); err != nil {
		return model.OperatorNote{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if err := validation.Text(input.Body, "body", 2, 2000); err != nil {
		return model.OperatorNote{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if input.ObservationID == "" && input.AlertID == "" {
		return model.OperatorNote{}, fmt.Errorf("%w: note target is required", ErrInvalidInput)
	}
	note := model.OperatorNote{ID: l.nextID("note"), ObservationID: input.ObservationID, AlertID: input.AlertID, Operator: input.Operator, Body: input.Body, CreatedAt: l.Now()}
	if err := l.store.SaveNote(ctx, note); err != nil {
		return model.OperatorNote{}, err
	}
	return note, nil
}

func (l *Lab) ListNotes(ctx context.Context, observationID, alertID string) ([]model.OperatorNote, error) {
	return l.store.ListNotes(ctx, observationID, alertID)
}

func (l *Lab) AddObservationNote(ctx context.Context, observationID, operator, body string) (model.OperatorNote, error) {
	return l.AddNote(ctx, model.NoteInput{ObservationID: observationID, Operator: operator, Body: body})
}
