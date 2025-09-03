package model

import "time"

// StatsEntry represents a fizzbuzz request statistics record in the database
type StatsEntry struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Int1      int       `gorm:"not null;uniqueIndex:idx_params" json:"-"`
	Int2      int       `gorm:"not null;uniqueIndex:idx_params" json:"-"`
	Limit     int       `gorm:"not null;uniqueIndex:idx_params" json:"-"`
	Str1      string    `gorm:"not null;size:100;uniqueIndex:idx_params" json:"-"`
	Str2      string    `gorm:"not null;size:100;uniqueIndex:idx_params" json:"-"`
	HitCount  int64     `gorm:"not null;default:0" json:"hit_count"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// StatsResponse represents the API response for the most frequent request
type StatsResponse struct {
	Request  FizzBuzzRequest `json:"request"`
	HitCount int64           `json:"hit_count"`
}