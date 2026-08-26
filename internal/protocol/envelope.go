package protocol

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/validation"
)

type Envelope struct {
	Version   int               `json:"version"`
	MessageID string            `json:"message_id"`
	SentAt    time.Time         `json:"sent_at"`
	Source    string            `json:"source"`
	Frames    []FrameRecord     `json:"frames"`
	Metadata  map[string]string `json:"metadata"`
}

type FrameRecord struct {
	ID            string  `json:"id"`
	ObservationID string  `json:"observation_id"`
	AntennaID     string  `json:"antenna_id"`
	Sequence      int     `json:"sequence"`
	CapturedAt    string  `json:"captured_at"`
	Signal        float64 `json:"signal"`
	Drift         float64 `json:"drift"`
	Completeness  float64 `json:"completeness"`
	PayloadHash   string  `json:"payload_hash"`
}

func (e Envelope) Valid() error {
	if e.Version != 1 {
		return fmt.Errorf("unsupported envelope version %d", e.Version)
	}
	if err := validation.Required(e.MessageID, "message_id"); err != nil {
		return err
	}
	if err := validation.Required(e.Source, "source"); err != nil {
		return err
	}
	if e.SentAt.IsZero() {
		return fmt.Errorf("sent_at is required")
	}
	if len(e.Frames) == 0 {
		return fmt.Errorf("frames are required")
	}
	for index, frame := range e.Frames {
		if err := frame.Valid(); err != nil {
			return fmt.Errorf("frame %d: %w", index, err)
		}
	}
	return nil
}

func (f FrameRecord) Valid() error {
	for name, value := range map[string]string{"id": f.ID, "observation_id": f.ObservationID, "antenna_id": f.AntennaID, "payload_hash": f.PayloadHash} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	if f.Sequence < 0 || f.Signal < 0 || f.Drift < 0 || f.Completeness < 0 || f.Completeness > 1 {
		return fmt.Errorf("frame measurements are invalid")
	}
	if _, err := time.Parse(time.RFC3339Nano, f.CapturedAt); err != nil {
		return fmt.Errorf("captured_at: %w", err)
	}
	return nil
}

func (e Envelope) JSON() ([]byte, error) {
	if err := e.Valid(); err != nil {
		return nil, err
	}
	return json.Marshal(e)
}

func ParseEnvelope(raw []byte) (Envelope, error) {
	var envelope Envelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return envelope, fmt.Errorf("decode envelope: %w", err)
	}
	if err := envelope.Valid(); err != nil {
		return envelope, fmt.Errorf("invalid envelope: %w", err)
	}
	return envelope, nil
}

func (f FrameRecord) Model() (model.CalibrationFrame, error) {
	if err := f.Valid(); err != nil {
		return model.CalibrationFrame{}, err
	}
	captured, err := time.Parse(time.RFC3339Nano, f.CapturedAt)
	if err != nil {
		return model.CalibrationFrame{}, err
	}
	return model.CalibrationFrame{ID: f.ID, ObservationID: f.ObservationID, AntennaID: f.AntennaID, Sequence: f.Sequence, CapturedAt: captured, Signal: f.Signal, Drift: f.Drift, Completeness: f.Completeness, PayloadHash: f.PayloadHash, State: model.FrameReceived}, nil
}

func FromModel(frame model.CalibrationFrame) FrameRecord {
	return FrameRecord{ID: frame.ID, ObservationID: frame.ObservationID, AntennaID: frame.AntennaID, Sequence: frame.Sequence, CapturedAt: frame.CapturedAt.UTC().Format(time.RFC3339Nano), Signal: frame.Signal, Drift: frame.Drift, Completeness: frame.Completeness, PayloadHash: frame.PayloadHash}
}
