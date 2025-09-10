package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/julietteengel/fizzbuzz-api/internal/mocks"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

func TestFizzBuzzService_GenerateFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		request  model.FizzBuzzRequest
		expected *model.FizzBuzzResponse
		wantErr  bool
		errMsg   string
	}{
		{
			name: "classic_fizzbuzz",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: &model.FizzBuzzResponse{
				Result: []string{
					"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz",
				},
				Count: 15,
			},
			wantErr: false,
		},
		{
			name: "custom_parameters",
			request: model.FizzBuzzRequest{
				Int1:  2,
				Int2:  4,
				Limit: 8,
				Str1:  "foo",
				Str2:  "bar",
			},
			expected: &model.FizzBuzzResponse{
				Result: []string{
					"1", "foo", "3", "foobar", "5", "foo", "7", "foobar",
				},
				Count: 8,
			},
			wantErr: false,
		},
		{
			name: "limit_one",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 1,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: &model.FizzBuzzResponse{
				Result: []string{"1"},
				Count:  1,
			},
			wantErr: false,
		},
		{
			name: "first_multiple_int1",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  7,
				Limit: 3,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: &model.FizzBuzzResponse{
				Result: []string{"1", "2", "fizz"},
				Count:  3,
			},
			wantErr: false,
		},
		{
			name: "first_multiple_int2",
			request: model.FizzBuzzRequest{
				Int1:  7,
				Int2:  3,
				Limit: 3,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: &model.FizzBuzzResponse{
				Result: []string{"1", "2", "buzz"},
				Count:  3,
			},
			wantErr: false,
		},
		// Note: Validation test cases removed as validation now happens only at controller level
		// The service trusts that it receives valid data from the controller
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStatsRepo := mocks.NewMockIStatsRepository(t)
			
			if !tt.wantErr {
				mockStatsRepo.EXPECT().RecordRequest(mock.Anything, tt.request).Return(nil).Once()
			}

			service := NewFizzBuzzService(mockStatsRepo)
			
			result, err := service.GenerateFizzBuzz(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Result, result.Result)
				assert.Equal(t, tt.expected.Count, result.Count)
				
				// Wait a bit for the async goroutine to complete
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}

func TestFizzBuzzService_GenerateFizzBuzz_StatsError(t *testing.T) {
	mockStatsRepo := mocks.NewMockIStatsRepository(t)
	
	request := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 5,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	mockStatsRepo.EXPECT().RecordRequest(mock.Anything, request).Return(assert.AnError).Once()

	service := NewFizzBuzzService(mockStatsRepo)
	
	result, err := service.GenerateFizzBuzz(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []string{"1", "2", "fizz", "4", "buzz"}, result.Result)
	assert.Equal(t, 5, result.Count)
	
	// Wait a bit for the async goroutine to complete
	time.Sleep(10 * time.Millisecond)
}