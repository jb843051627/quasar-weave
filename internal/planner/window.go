package planner

import (
	"fmt"
	"sort"
	"time"
)

type Window struct {
	ID       string    `json:"id"`
	Target   string    `json:"target"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Priority int       `json:"priority"`
	Band     string    `json:"band"`
	Reserved bool      `json:"reserved"`
}

func (w Window) Duration() time.Duration { return w.End.Sub(w.Start) }

func (w Window) Valid() error {
	if w.ID == "" || w.Target == "" || w.Band == "" {
		return fmt.Errorf("window identity is incomplete")
	}
	if w.Start.IsZero() || w.End.IsZero() || !w.End.After(w.Start) {
		return fmt.Errorf("window interval is invalid")
	}
	if w.Priority < 0 || w.Priority > 100 {
		return fmt.Errorf("priority is outside range")
	}
	return nil
}

func SortWindows(windows []Window) []Window {
	result := append([]Window(nil), windows...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority > result[j].Priority
		}
		return result[i].Start.Before(result[j].Start)
	})
	return result
}

func Overlap(left, right Window) bool {
	return left.Start.Before(right.End) && right.Start.Before(left.End)
}

func Trim(w Window, start, end time.Time) (Window, error) {
	if start.After(w.Start) {
		w.Start = start
	}
	if end.Before(w.End) {
		w.End = end
	}
	if err := w.Valid(); err != nil {
		return Window{}, err
	}
	return w, nil
}
