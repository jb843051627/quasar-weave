package model

import "time"

type Antenna struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Station       string        `json:"station"`
	Band          string        `json:"band"`
	Status        AntennaStatus `json:"status"`
	Enabled       bool          `json:"enabled"`
	Latitude      float64       `json:"latitude"`
	Longitude     float64       `json:"longitude"`
	LastHeartbeat time.Time     `json:"last_heartbeat"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type AntennaInput struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Station   string  `json:"station"`
	Band      string  `json:"band"`
	Enabled   bool    `json:"enabled"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (a Antenna) Healthy(now time.Time, maxSilence time.Duration) bool {
	if !a.Enabled || a.Status != AntennaReady || a.LastHeartbeat.IsZero() {
		return false
	}
	return now.Sub(a.LastHeartbeat) <= maxSilence
}

func (a Antenna) Coordinates() [2]float64 {
	return [2]float64{a.Latitude, a.Longitude}
}
