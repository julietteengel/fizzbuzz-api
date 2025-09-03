package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/julietteengel/fizzbuzz-api/internal/mocks"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

func TestStatsService_GetMostFrequent(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockIStatsRepository)
		expectedResult *model.StatsResponse
		expectedError  bool
	}{
		{
			name: "successful_stats_retrieval",
			mockSetup: func(mockRepo *mocks.MockIStatsRepository) {
				expectedResponse := &model.StatsResponse{
					Request: model.FizzBuzzRequest{
						Int1:  3,
						Int2:  5,
						Limit: 100,
						Str1:  "fizz",
						Str2:  "buzz",
					},
					HitCount: 42,
				}
				mockRepo.EXPECT().GetMostFrequent(mock.Anything).Return(expectedResponse, nil).Once()
			},
			expectedResult: &model.StatsResponse{
				Request: model.FizzBuzzRequest{
					Int1:  3,
					Int2:  5,
					Limit: 100,
					Str1:  "fizz",
					Str2:  "buzz",
				},
				HitCount: 42,
			},
			expectedError: false,
		},
		{
			name: "no_stats_available",
			mockSetup: func(mockRepo *mocks.MockIStatsRepository) {
				mockRepo.EXPECT().GetMostFrequent(mock.Anything).Return(nil, nil).Once()
			},
			expectedResult: nil,
			expectedError:  false,
		},
		{
			name: "repository_error",
			mockSetup: func(mockRepo *mocks.MockIStatsRepository) {
				mockRepo.EXPECT().GetMostFrequent(mock.Anything).Return(nil, assert.AnError).Once()
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockIStatsRepository(t)
			tt.mockSetup(mockRepo)

			service := NewStatsService(mockRepo)

			result, err := service.GetMostFrequent(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectedResult == nil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, tt.expectedResult.Request, result.Request)
					assert.Equal(t, tt.expectedResult.HitCount, result.HitCount)
				}
			}
		})
	}
}

func TestStatsService_GetMostFrequent_ContextPassing(t *testing.T) {
	mockRepo := mocks.NewMockIStatsRepository(t)
	
	ctx := context.WithValue(context.Background(), "test_key", "test_value")
	
	expectedResponse := &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  2,
			Int2:  7,
			Limit: 50,
			Str1:  "foo",
			Str2:  "bar",
		},
		HitCount: 15,
	}

	mockRepo.EXPECT().GetMostFrequent(ctx).Return(expectedResponse, nil).Once()

	service := NewStatsService(mockRepo)
	
	result, err := service.GetMostFrequent(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse, result)
}