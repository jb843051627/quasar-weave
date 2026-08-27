package model

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("record not found")
	ErrInvalidState   = errors.New("invalid state transition")
	ErrAlreadyExists  = errors.New("record already exists")
	ErrQualityFailure = errors.New("quality gate rejected frame")
	ErrQueueClosed    = errors.New("processing queue is closed")
)

type AntennaStatus string

const (
	AntennaReady   AntennaStatus = "ready"
	AntennaMuted   AntennaStatus = "muted"
	AntennaOffline AntennaStatus = "offline"
	AntennaFault   AntennaStatus = "fault"
)

func (s AntennaStatus) Valid() bool {
	switch s {
	case AntennaReady, AntennaMuted, AntennaOffline, AntennaFault:
		return true
	default:
		return false
	}
}

type ObservationStatus string

const (
	ObservationPlanned     ObservationStatus = "planned"
	ObservationCapturing   ObservationStatus = "capturing"
	ObservationCalibrating ObservationStatus = "calibrating"
	ObservationQualified   ObservationStatus = "qualified"
	ObservationFailed      ObservationStatus = "failed"
	ObservationArchived    ObservationStatus = "archived"
)

func (s ObservationStatus) Terminal() bool {
	return s == ObservationArchived
}

func (s ObservationStatus) Valid() bool {
	switch s {
	case ObservationPlanned, ObservationCapturing, ObservationCalibrating,
		ObservationQualified, ObservationFailed, ObservationArchived:
		return true
	default:
		return false
	}
}

type FrameState string

const (
	FrameReceived FrameState = "received"
	FrameQueued   FrameState = "queued"
	FrameChecked  FrameState = "checked"
	FrameRejected FrameState = "rejected"
)

type RetryState string

const (
	RetryPending  RetryState = "pending"
	RetryRunning  RetryState = "running"
	RetryFinished RetryState = "finished"
	RetryCanceled RetryState = "canceled"
)

type AlertState string

const (
	AlertOpen     AlertState = "open"
	AlertAcked    AlertState = "acked"
	AlertResolved AlertState = "resolved"
)

func ValidateID(kind, id string) error {
	if id == "" {
		return fmt.Errorf("%s id is required", kind)
	}
	if len(id) > 80 {
		return fmt.Errorf("%s id is too long", kind)
	}
	return nil
}

func CanTransition(from, to ObservationStatus) bool {
	if !from.Valid() || !to.Valid() || from == to {
		return false
	}
	switch from {
	case ObservationPlanned:
		return to == ObservationCapturing || to == ObservationFailed
	case ObservationCapturing:
		return to == ObservationCalibrating || to == ObservationFailed
	case ObservationCalibrating:
		return to == ObservationQualified || to == ObservationFailed
	case ObservationQualified:
		return to == ObservationArchived || to == ObservationFailed
	case ObservationFailed:
		return to == ObservationCapturing || to == ObservationArchived
	default:
		return false
	}
}
