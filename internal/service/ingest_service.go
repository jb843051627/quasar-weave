package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/protocol"
)

func (l *Lab) IngestEnvelope(ctx context.Context, envelope protocol.Envelope) (int, error) {
	if err := envelope.Valid(); err != nil {
		return 0, fmt.Errorf("validate ingest envelope: %w", err)
	}
	batch := protocol.NewBatch(envelope.Source)
	if err := batch.Add(envelope); err != nil {
		return 0, err
	}
	count := 0
	for _, input := range batch.Ordered() {
		if err := ctx.Err(); err != nil {
			return count, err
		}
		if _, done, err := l.SubmitFrame(ctx, input); err != nil {
			return count, err
		} else if done != nil {
			select {
			case err := <-done:
				if err != nil {
					return count, err
				}
			case <-ctx.Done():
				return count, ctx.Err()
			}
		}
		count++
	}
	return count, nil
}

func (l *Lab) IngestStream(ctx context.Context, reader io.Reader, source string, now time.Time) (int, error) {
	stream := protocol.NewStream(0, 10*time.Minute)
	records, err := stream.Read(ctx, reader, now)
	if err != nil {
		return 0, err
	}
	envelope := protocol.Envelope{Version: 1, MessageID: l.nextID("message"), SentAt: now, Source: source, Frames: records}
	return l.IngestEnvelope(ctx, envelope)
}

func frameInput(frame model.CalibrationFrame) model.FrameInput {
	return model.FrameInput{ID: frame.ID, ObservationID: frame.ObservationID, AntennaID: frame.AntennaID, Sequence: frame.Sequence, CapturedAt: frame.CapturedAt, Signal: frame.Signal, Drift: frame.Drift, Completeness: frame.Completeness, PayloadHash: frame.PayloadHash}
}
