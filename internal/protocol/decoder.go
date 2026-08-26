package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type LineDecoder struct {
	maxLine int
}

func NewLineDecoder(maxLine int) *LineDecoder {
	if maxLine < 256 {
		maxLine = 64 * 1024
	}
	return &LineDecoder{maxLine: maxLine}
}

func (d *LineDecoder) Decode(reader io.Reader) ([]FrameRecord, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), d.maxLine)
	frames := make([]FrameRecord, 0)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		frame, err := ParseLine(text)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		frames = append(frames, frame)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan frame stream: %v", err)
	}
	return frames, nil
}

func ParseLine(line string) (FrameRecord, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 9 {
		return FrameRecord{}, fmt.Errorf("expected 9 fields, got %d", len(parts))
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	sequence, err := strconv.Atoi(parts[3])
	if err != nil {
		return FrameRecord{}, fmt.Errorf("sequence: %w", err)
	}
	signal, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return FrameRecord{}, fmt.Errorf("signal: %w", err)
	}
	drift, err := strconv.ParseFloat(parts[6], 64)
	if err != nil {
		return FrameRecord{}, fmt.Errorf("drift: %w", err)
	}
	completeness, err := strconv.ParseFloat(parts[7], 64)
	if err != nil {
		return FrameRecord{}, fmt.Errorf("completeness: %w", err)
	}
	frame := FrameRecord{ID: parts[0], ObservationID: parts[1], AntennaID: parts[2], Sequence: sequence, CapturedAt: parts[4], Signal: signal, Drift: drift, Completeness: completeness, PayloadHash: parts[8]}
	if err := frame.Valid(); err != nil {
		return FrameRecord{}, err
	}
	return frame, nil
}

func EncodeLine(frame FrameRecord) (string, error) {
	if err := frame.Valid(); err != nil {
		return "", err
	}
	return strings.Join([]string{frame.ID, frame.ObservationID, frame.AntennaID, strconv.Itoa(frame.Sequence), frame.CapturedAt, strconv.FormatFloat(frame.Signal, 'f', 6, 64), strconv.FormatFloat(frame.Drift, 'f', 6, 64), strconv.FormatFloat(frame.Completeness, 'f', 6, 64), frame.PayloadHash}, "|"), nil
}

func ParseLines(lines []string) ([]FrameRecord, error) {
	return NewLineDecoder(0).Decode(strings.NewReader(strings.Join(lines, "\n")))
}

func EncodeLines(frames []FrameRecord) ([]byte, error) {
	var buffer bytes.Buffer
	for _, frame := range frames {
		line, err := EncodeLine(frame)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(line)
		buffer.WriteByte('\n')
	}
	return buffer.Bytes(), nil
}

func NormalizeCapturedAt(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }
