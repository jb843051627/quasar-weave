package planner

import (
	"fmt"
	"strings"
	"time"
)

type Strategy struct {
	Name             string        `json:"name"`
	Band             string        `json:"band"`
	FrameInterval    time.Duration `json:"frame_interval"`
	ExpectedFrames   int           `json:"expected_frames"`
	RetryLimit       int           `json:"retry_limit"`
	RequireHeartbeat bool          `json:"require_heartbeat"`
}

func (s Strategy) Valid() error {
	if strings.TrimSpace(s.Name) == "" || strings.TrimSpace(s.Band) == "" {
		return fmt.Errorf("strategy name and band are required")
	}
	if s.FrameInterval <= 0 || s.ExpectedFrames <= 0 || s.RetryLimit < 0 {
		return fmt.Errorf("strategy timing values are invalid")
	}
	return nil
}

func (s Strategy) CaptureDuration() time.Duration {
	return time.Duration(s.ExpectedFrames) * s.FrameInterval
}

func (s Strategy) DueFrames(start, now time.Time) int {
	if now.Before(start) || s.FrameInterval <= 0 {
		return 0
	}
	count := int(now.Sub(start) / s.FrameInterval)
	if count > s.ExpectedFrames {
		return s.ExpectedFrames
	}
	return count
}

func (s Strategy) WithExpectedFrames(count int) Strategy {
	if count > 0 {
		s.ExpectedFrames = count
	}
	return s
}
