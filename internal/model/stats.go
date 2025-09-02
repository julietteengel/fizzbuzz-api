package model

import "time"

type StatsEntry struct {
	Parameters FizzBuzzRequest `json:"parameters"`
	HitCount   int64           `json:"hit_count"`
	LastAccess time.Time       `json:"last_access"`
	FirstSeen  time.Time       `json:"first_seen"`
}

type StatsResponse struct {
	MostFrequent *StatsEntry `json:"most_frequent,omitempty"`
	Message      string      `json:"message,omitempty"`
}

type AllStatsResponse struct {
	Stats      []StatsEntry `json:"stats"`
	TotalCount int          `json:"total_count"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime,omitempty"`
}

type RequestMetadata struct {
	RequestID  string    `json:"request_id"`
	ClientIP   string    `json:"client_ip"`
	UserAgent  string    `json:"user_agent"`
	ReceivedAt time.Time `json:"received_at"`
}