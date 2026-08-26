package validation

import (
	"fmt"
	"regexp"
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{2,79}$`)

func ID(value string) error {
	if !idPattern.MatchString(value) {
		return fmt.Errorf("id must contain lowercase letters, digits, '_' or '-' and be 3-80 chars")
	}
	return nil
}

func Required(value, field string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func OneOf(value, field string, options ...string) error {
	for _, option := range options {
		if value == option {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %v", field, options)
}
