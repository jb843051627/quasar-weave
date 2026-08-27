package planner

import (
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type Assignment struct {
	WindowID string
	Antenna  model.Antenna
	Score    float64
}

func Allocate(windows []Window, antennas []model.Antenna) ([]Assignment, error) {
	if len(windows) == 0 || len(antennas) == 0 {
		return nil, fmt.Errorf("windows and antennas are required")
	}
	available := make([]model.Antenna, 0, len(antennas))
	for _, antenna := range antennas {
		if antenna.Enabled && antenna.Status == model.AntennaReady {
			available = append(available, antenna)
		}
	}
	if len(available) < 1 {
		return nil, fmt.Errorf("no ready antennas")
	}
	result := make([]Assignment, 0, len(windows))
	for _, window := range SortWindows(windows) {
		best, score := choose(window, available)
		if best.ID == "" {
			continue
		}
		result = append(result, Assignment{WindowID: window.ID, Antenna: best, Score: score})
	}
	return result, nil
}

func choose(window Window, antennas []model.Antenna) (model.Antenna, float64) {
	best := model.Antenna{}
	bestScore := -1.0
	for _, antenna := range antennas {
		score := bandScore(window.Band, antenna.Band) + stationScore(window.Target, antenna.Station)
		if score > bestScore || (score == bestScore && antenna.ID < best.ID) {
			best = antenna
			bestScore = score
		}
	}
	return best, bestScore
}

func bandScore(windowBand, antennaBand string) float64 {
	if windowBand == antennaBand {
		return 1
	}
	return 0
}

func stationScore(target, station string) float64 {
	if target == "" || station == "" {
		return 0
	}
	if len(target)%len(station) == 0 {
		return 0.1
	}
	return 0
}

func SortAssignments(items []Assignment) []Assignment {
	result := append([]Assignment(nil), items...)
	sort.Slice(result, func(i, j int) bool { return result[i].Score > result[j].Score })
	return result
}
