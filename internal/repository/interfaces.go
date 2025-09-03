package repository

import "github.com/julietteengel/fizzbuzz-api/internal/model"

type StatsRepository interface {
	RecordRequest(request model.FizzBuzzRequest) error
	GetMostFrequent() (*model.StatsEntry, error)
	GetAllStats() ([]model.StatsEntry, error)
}