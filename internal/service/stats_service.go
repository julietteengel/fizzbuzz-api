package service

import (
	"context"

	"github.com/julietteengel/fizzbuzz-api/internal/model"
	"github.com/julietteengel/fizzbuzz-api/internal/repository"
)

type IStatsService interface {
	GetMostFrequent(ctx context.Context) (*model.StatsResponse, error)
}

type statsService struct {
	statsRepo repository.IStatsRepository
}

func NewStatsService(statsRepo repository.IStatsRepository) IStatsService {
	return &statsService{
		statsRepo: statsRepo,
	}
}

func (s *statsService) GetMostFrequent(ctx context.Context) (*model.StatsResponse, error) {
	return s.statsRepo.GetMostFrequent(ctx)
}