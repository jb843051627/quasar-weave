package clock

import (
	"fmt"
	"time"
)

const DefaultLayout = time.RFC3339Nano

func Parse(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("time value is empty")
	}
	parsed, err := time.Parse(DefaultLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp: %w", err)
	}
	return parsed.UTC(), nil
}

func Format(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(DefaultLayout)
}

func Window(now time.Time, before, after time.Duration) (time.Time, time.Time) {
	return now.Add(-before), now.Add(after)
}
