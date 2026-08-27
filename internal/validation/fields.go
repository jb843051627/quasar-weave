package validation

import (
	"fmt"
	"strings"
)

func Text(value, field string, min, max int) error {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < min {
		return fmt.Errorf("%s is too short", field)
	}
	if len(trimmed) > max {
		return fmt.Errorf("%s is too long", field)
	}
	return nil
}

func Positive(value int, field string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", field)
	}
	return nil
}

func Range(value, min, max float64, field string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %.3f and %.3f", field, min, max)
	}
	return nil
}
