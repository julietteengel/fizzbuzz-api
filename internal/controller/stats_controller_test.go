package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/julietteengel/fizzbuzz-api/internal/mocks"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

func TestStatsController_GetStats_Success(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIStatsService(t)
	controller := NewStatsController(mockService)

	expectedStats := &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  3,
			Int2:  5,
			Limit: 100,
			Str1:  "fizz",
			Str2:  "buzz",
		},
		HitCount: 42,
	}

	mockService.EXPECT().GetMostFrequent(mock.Anything).Return(expectedStats, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetStats(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.StatsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStats.Request, response.Request)
	assert.Equal(t, expectedStats.HitCount, response.HitCount)
}

func TestStatsController_GetStats_NoStats(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIStatsService(t)
	controller := NewStatsController(mockService)

	mockService.EXPECT().GetMostFrequent(mock.Anything).Return(nil, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetStats(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestStatsController_GetStats_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIStatsService(t)
	controller := NewStatsController(mockService)

	mockService.EXPECT().GetMostFrequent(mock.Anything).Return(nil, assert.AnError).Once()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetStats(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
}

func TestStatsController_GetStats_ContextPassing(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIStatsService(t)
	controller := NewStatsController(mockService)

	expectedStats := &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  2,
			Int2:  7,
			Limit: 50,
			Str1:  "foo",
			Str2:  "bar",
		},
		HitCount: 15,
	}

	// We expect the service to be called with any context
	mockService.EXPECT().GetMostFrequent(mock.Anything).Return(expectedStats, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetStats(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.StatsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStats, &response)
}

func TestStatsController_GetStats_ZeroHitCount(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIStatsService(t)
	controller := NewStatsController(mockService)

	// Edge case: stats exist but hit count is 0
	statsWithZeroHits := &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  1,
			Int2:  1,
			Limit: 1,
			Str1:  "test",
			Str2:  "test",
		},
		HitCount: 0,
	}

	mockService.EXPECT().GetMostFrequent(mock.Anything).Return(statsWithZeroHits, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetStats(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.StatsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, statsWithZeroHits.Request, response.Request)
	assert.Equal(t, int64(0), response.HitCount)
}