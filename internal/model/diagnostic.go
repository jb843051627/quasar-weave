package model

import (
	"sort"
	"strings"
	"time"
)

type DiagnosticKind string

const (
	DiagnosticSignal    DiagnosticKind = "signal"
	DiagnosticDrift     DiagnosticKind = "drift"
	DiagnosticCoverage  DiagnosticKind = "coverage"
	DiagnosticHeartbeat DiagnosticKind = "heartbeat"
	DiagnosticStorage   DiagnosticKind = "storage"
)

type Diagnostic struct {
	ID            string         `json:"id"`
	ObservationID string         `json:"observation_id"`
	AntennaID     string         `json:"antenna_id"`
	Kind          DiagnosticKind `json:"kind"`
	Severity      string         `json:"severity"`
	Message       string         `json:"message"`
	Evidence      []string       `json:"evidence"`
	CreatedAt     time.Time      `json:"created_at"`
}

type DiagnosticSet struct {
	Items []Diagnostic `json:"items"`
}

func (d Diagnostic) Valid() bool {
	if d.ID == "" || d.Kind == "" || d.Message == "" || d.CreatedAt.IsZero() {
		return false
	}
	switch d.Severity {
	case "info", "warning", "critical":
		return true
	default:
		return false
	}
}

func (d Diagnostic) SearchText() string {
	parts := append([]string{d.Kind.String(), d.Severity, d.Message}, d.Evidence...)
	return strings.ToLower(strings.Join(parts, " "))
}

func (k DiagnosticKind) String() string { return string(k) }

func (s *DiagnosticSet) Add(item Diagnostic) {
	if item.Valid() {
		s.Items = append(s.Items, item)
	}
}

func (s DiagnosticSet) BySeverity(severity string) []Diagnostic {
	result := make([]Diagnostic, 0)
	for _, item := range s.Items {
		if item.Severity == severity {
			result = append(result, item)
		}
	}
	return result
}

func (s DiagnosticSet) Ordered() []Diagnostic {
	result := append([]Diagnostic(nil), s.Items...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Severity != result[j].Severity {
			return severityRank(result[i].Severity) > severityRank(result[j].Severity)
		}
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

func severityRank(value string) int {
	switch value {
	case "critical":
		return 3
	case "warning":
		return 2
	default:
		return 1
	}
}

func (s DiagnosticSet) HasCritical() bool {
	for _, item := range s.Items {
		if item.Severity == "critical" {
			return true
		}
	}
	return false
}
