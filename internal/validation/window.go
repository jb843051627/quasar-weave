package validation

import (
	"fmt"
	"time"
)

func InWindow(value, start, end time.Time) bool {
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return false
	}
	return !value.Before(start) && !value.After(end)
}

func RequireRecent(value, now time.Time, maxAge time.Duration) error {
	if value.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	if value.After(now) {
		return fmt.Errorf("timestamp is in the future")
	}
	if now.Sub(value) > maxAge {
		return fmt.Errorf("timestamp is stale")
	}
	return nil
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
