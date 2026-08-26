package protocol

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Stream struct {
	decoder *LineDecoder
	maxAge  time.Duration
}

func NewStream(maxLine int, maxAge time.Duration) *Stream {
	if maxAge <= 0 {
		maxAge = 10 * time.Minute
	}
	return &Stream{decoder: NewLineDecoder(maxLine), maxAge: maxAge}
}

func (s *Stream) Read(ctx context.Context, reader io.Reader, now time.Time) ([]FrameRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if reader == nil {
		return nil, fmt.Errorf("stream reader is required")
	}
	frames, err := s.decoder.Decode(reader)
	if err != nil {
		return nil, err
	}
	result := make([]FrameRecord, 0, len(frames))
	for _, frame := range frames {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		captured, parseErr := time.Parse(time.RFC3339Nano, frame.CapturedAt)
		if parseErr == nil && now.Sub(captured) <= s.maxAge {
			result = append(result, frame)
		}
	}
	return result, nil
}
