package protocol

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func FormatFrameLine(envelope Envelope) (string, error) {
	if err := envelope.Valid(); err != nil {
		return "", err
	}
	if len(envelope.Frames) == 0 {
		return "", fmt.Errorf("frame list is empty")
	}
	f := envelope.Frames[0]
	relation := f.ObservationID + ":" + f.AntennaID
	values := []string{envelope.MessageID, envelope.Source, f.ID, strconv.Itoa(f.Sequence), envelope.SentAt.Format(time.RFC3339Nano), f.CapturedAt, strconv.FormatFloat(f.Signal, 'g', -1, 64), strconv.FormatFloat(f.Drift, 'g', -1, 64), strconv.FormatFloat(f.Completeness, 'g', -1, 64), relation}
	return strings.Join(values, "|"), nil
}

func ParseFrameLine(line string, maxLine int) (model.FrameInput, error) {
	decoder := NewLineDecoder(maxLine)
	envelopes, err := decoder.Decode(strings.NewReader(line + "\n"))
	if err != nil {
		return model.FrameInput{}, err
	}
	if len(envelopes) != 1 {
		return model.FrameInput{}, fmt.Errorf("expected one frame")
	}
	frame, err := envelopes[0].Model()
	if err != nil {
		return model.FrameInput{}, err
	}
	return model.FrameInput{ID: frame.ID, ObservationID: frame.ObservationID, AntennaID: frame.AntennaID, Sequence: frame.Sequence, CapturedAt: frame.CapturedAt, Signal: frame.Signal, Drift: frame.Drift, Completeness: frame.Completeness, PayloadHash: frame.PayloadHash}, nil
}
