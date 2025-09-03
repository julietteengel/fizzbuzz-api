package model

import "time"

// HealthCheckResponse represents the health check endpoint response
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime,omitempty"`
}