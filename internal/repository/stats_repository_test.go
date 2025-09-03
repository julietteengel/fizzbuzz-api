package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

func TestStatsRepository_Memory_RecordRequest(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg)

	request1 := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	// First record
	err := repo.RecordRequest(context.Background(), request1)
	assert.NoError(t, err)

	// Second record of the same request
	err = repo.RecordRequest(context.Background(), request1)
	assert.NoError(t, err)

	// Get most frequent
	result, err := repo.GetMostFrequent(context.Background())
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, request1, result.Request)
	assert.Equal(t, int64(2), result.HitCount)
}

func TestStatsRepository_Memory_GetMostFrequent_Multiple(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg)

	request1 := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	request2 := model.FizzBuzzRequest{
		Int1:  2,
		Int2:  7,
		Limit: 50,
		Str1:  "foo",
		Str2:  "bar",
	}

	// Record request1 three times
	for i := 0; i < 3; i++ {
		err := repo.RecordRequest(context.Background(), request1)
		assert.NoError(t, err)
	}

	// Record request2 five times (should become most frequent)
	for i := 0; i < 5; i++ {
		err := repo.RecordRequest(context.Background(), request2)
		assert.NoError(t, err)
	}

	// Get most frequent - should be request2
	result, err := repo.GetMostFrequent(context.Background())
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, request2, result.Request)
	assert.Equal(t, int64(5), result.HitCount)
}

func TestStatsRepository_Memory_GetMostFrequent_Empty(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg)

	result, err := repo.GetMostFrequent(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestStatsRepository_Memory_ThreadSafety(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg)

	request := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	// Run concurrent operations to test thread safety
	numGoroutines := 10
	numRecords := 5

	done := make(chan bool, numGoroutines)

	// Start multiple goroutines recording the same request
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numRecords; j++ {
				err := repo.RecordRequest(context.Background(), request)
				assert.NoError(t, err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check final count
	result, err := repo.GetMostFrequent(context.Background())
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, request, result.Request)
	assert.Equal(t, int64(numGoroutines*numRecords), result.HitCount)
}

func TestStatsRepository_Memory_GenerateKey(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg).(*statsRepository)

	request1 := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	request2 := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	request3 := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 101, // Different limit
		Str1:  "fizz",
		Str2:  "buzz",
	}

	key1 := repo.generateKey(request1)
	key2 := repo.generateKey(request2)
	key3 := repo.generateKey(request3)

	// Same requests should generate same keys
	assert.Equal(t, key1, key2)
	// Different requests should generate different keys
	assert.NotEqual(t, key1, key3)
}

func TestStatsRepository_Database_Mode(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "postgres",
		},
	}

	// Create repository with nil database (we're not testing actual DB operations here)
	// This tests the initialization and mode selection
	repo := NewStatsRepository(nil, cfg).(*statsRepository)

	assert.False(t, repo.useMemory)
	assert.Nil(t, repo.db) // We passed nil, so it should be nil
	assert.NotNil(t, repo.memStats) // Should still be initialized
}

func TestStatsRepository_Memory_MultipleDistinctRequests(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			StatsStorage: "memory",
		},
	}
	repo := NewStatsRepository(nil, cfg)

	requests := []model.FizzBuzzRequest{
		{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
		{Int1: 2, Int2: 7, Limit: 50, Str1: "foo", Str2: "bar"},
		{Int1: 4, Int2: 6, Limit: 75, Str1: "ping", Str2: "pong"},
		{Int1: 1, Int2: 10, Limit: 200, Str1: "a", Str2: "b"},
	}

	hitCounts := []int{1, 5, 3, 2}

	// Record different requests different number of times
	for i, request := range requests {
		for j := 0; j < hitCounts[i]; j++ {
			err := repo.RecordRequest(context.Background(), request)
			assert.NoError(t, err)
		}
	}

	// Get most frequent - should be requests[1] with 5 hits
	result, err := repo.GetMostFrequent(context.Background())
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, requests[1], result.Request)
	assert.Equal(t, int64(5), result.HitCount)
}