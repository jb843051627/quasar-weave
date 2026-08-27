package validation

import (
	"fmt"
	"math"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func TelemetryPoint(point model.TelemetryPoint, now time.Time, maxAge time.Duration) error {
	if !point.Finite() {
		return fmt.Errorf("telemetry point is incomplete or non-finite")
	}
	if point.CapturedAt.After(now) {
		return fmt.Errorf("telemetry point is in the future")
	}
	if maxAge > 0 && now.Sub(point.CapturedAt) > maxAge {
		return fmt.Errorf("telemetry point is stale")
	}
	return nil
}

func TelemetryValue(value, minimum, maximum float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("telemetry value is not finite")
	}
	if value < minimum || value > maximum {
		return fmt.Errorf("telemetry value is outside range")
	}
	return nil
}

func Window(start, end time.Time, minimum time.Duration) error {
	if start.IsZero() || end.IsZero() || !end.After(start) {
		return fmt.Errorf("time window is invalid")
	}
	if end.Sub(start) < minimum {
		return fmt.Errorf("time window is too short")
	}
	return nil
}
