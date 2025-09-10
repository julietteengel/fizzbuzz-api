package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

type IStatsRepository interface {
	RecordRequest(ctx context.Context, request model.FizzBuzzRequest) error
	GetMostFrequent(ctx context.Context) (*model.StatsResponse, error)
}

type statsRepository struct {
	db        *gorm.DB
	memStats  map[string]*model.StatsEntry
	memMutex  sync.RWMutex
	useMemory bool
}

func NewStatsRepository(database *gorm.DB, cfg *config.Config) IStatsRepository {
	useMemory := cfg.Database.StatsStorage == "memory"
	return &statsRepository{
		db:        database,
		memStats:  make(map[string]*model.StatsEntry),
		useMemory: useMemory,
	}
}


func (r *statsRepository) RecordRequest(ctx context.Context, request model.FizzBuzzRequest) error {
	if r.useMemory {
		return r.recordInMemory(request)
	}
	return r.recordInDatabase(ctx, request)
}

func (r *statsRepository) GetMostFrequent(ctx context.Context) (*model.StatsResponse, error) {
	if r.useMemory {
		return r.getMostFrequentFromMemory()
	}
	return r.getMostFrequentFromDatabase(ctx)
}

func (r *statsRepository) recordInMemory(request model.FizzBuzzRequest) error {
	r.memMutex.Lock()
	defer r.memMutex.Unlock()

	key := r.generateKey(request)
	if entry, exists := r.memStats[key]; exists {
		entry.HitCount++
		entry.UpdatedAt = time.Now()
	} else {
		r.memStats[key] = &model.StatsEntry{
			Int1:      request.Int1,
			Int2:      request.Int2,
			Limit:     request.Limit,
			Str1:      request.Str1,
			Str2:      request.Str2,
			HitCount:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return nil
}

func (r *statsRepository) getMostFrequentFromMemory() (*model.StatsResponse, error) {
	r.memMutex.RLock()
	defer r.memMutex.RUnlock()

	var mostFrequent *model.StatsEntry
	for _, entry := range r.memStats {
		if mostFrequent == nil || entry.HitCount > mostFrequent.HitCount {
			mostFrequent = entry
		}
	}

	if mostFrequent == nil {
		return nil, nil
	}

	return &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  mostFrequent.Int1,
			Int2:  mostFrequent.Int2,
			Limit: mostFrequent.Limit,
			Str1:  mostFrequent.Str1,
			Str2:  mostFrequent.Str2,
		},
		HitCount: mostFrequent.HitCount,
	}, nil
}

func (r *statsRepository) recordInDatabase(ctx context.Context, request model.FizzBuzzRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		entry := model.StatsEntry{
			Int1:  request.Int1,
			Int2:  request.Int2,
			Limit: request.Limit,
			Str1:  request.Str1,
			Str2:  request.Str2,
		}

		result := tx.Where(&entry).First(&entry)
		if result.Error == gorm.ErrRecordNotFound {
			entry.HitCount = 1
			return tx.Create(&entry).Error
		}
		if result.Error != nil {
			return result.Error
		}

		entry.HitCount++
		return tx.Save(&entry).Error
	})
}

func (r *statsRepository) getMostFrequentFromDatabase(ctx context.Context) (*model.StatsResponse, error) {
	var entry model.StatsEntry
	result := r.db.WithContext(ctx).Order("hit_count DESC").First(&entry)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  entry.Int1,
			Int2:  entry.Int2,
			Limit: entry.Limit,
			Str1:  entry.Str1,
			Str2:  entry.Str2,
		},
		HitCount: entry.HitCount,
	}, nil
}

func (r *statsRepository) generateKey(request model.FizzBuzzRequest) string {
	return fmt.Sprintf("%d_%d_%d_%s_%s", request.Int1, request.Int2, request.Limit, request.Str1, request.Str2)
}