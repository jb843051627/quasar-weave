package protocol

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Batch struct {
	Source     string
	Frames     []model.FrameInput
	Messages   []string
	duplicates int
	mu         sync.Mutex
}

func NewBatch(source string) *Batch {
	return &Batch{Source: source, Frames: make([]model.FrameInput, 0), Messages: make([]string, 0)}
}

func (b *Batch) Add(envelope Envelope) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Source != "" && envelope.Source != b.Source {
		return fmt.Errorf("source mismatch: %s != %s", envelope.Source, b.Source)
	}
	for _, record := range envelope.Frames {
		frame, err := record.Model()
		if err != nil {
			return err
		}
		input := model.FrameInput{ID: frame.ID, ObservationID: frame.ObservationID, AntennaID: frame.AntennaID, Sequence: frame.Sequence, CapturedAt: frame.CapturedAt, Signal: frame.Signal, Drift: frame.Drift, Completeness: frame.Completeness, PayloadHash: frame.PayloadHash}
		duplicate := false
		for _, existing := range b.Frames {
			if existing.ID == input.ID {
				duplicate = true
				break
			}
		}
		if duplicate {
			b.duplicates++
			continue
		}
		b.Frames = append(b.Frames, input)
	}
	b.Messages = append(b.Messages, envelope.MessageID)
	return nil
}

func (b *Batch) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.Frames)
}

func (b *Batch) Duplicates() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.duplicates
}

func (b *Batch) Ordered() []model.FrameInput {
	b.mu.Lock()
	defer b.mu.Unlock()
	frames := append([]model.FrameInput(nil), b.Frames...)
	sort.SliceStable(frames, func(i, j int) bool {
		if frames[i].ObservationID != frames[j].ObservationID {
			return frames[i].ObservationID < frames[j].ObservationID
		}
		return frames[i].Sequence < frames[j].Sequence
	})
	return frames
}

func (b *Batch) Drain(ctx context.Context, submit func(context.Context, model.FrameInput) error) error {
	if submit == nil {
		return fmt.Errorf("submit function is required")
	}
	for _, frame := range b.Ordered() {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := submit(ctx, frame); err != nil {
			return fmt.Errorf("submit frame %s: %w", frame.ID, err)
		}
	}
	return nil
}
