package controller

import (
	"bytes"
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

func TestFizzBuzzController_GenerateFizzBuzz_Success(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIFizzBuzzService(t)
	controller := NewFizzBuzzController(mockService)

	request := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	expectedResponse := &model.FizzBuzzResponse{
		Result: []string{
			"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz",
		},
		Count: 15,
	}

	mockService.EXPECT().GenerateFizzBuzz(mock.Anything, request).Return(expectedResponse, nil).Once()

	// Marshal request to JSON
	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/fizzbuzz", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GenerateFizzBuzz(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.FizzBuzzResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Result, response.Result)
	assert.Equal(t, expectedResponse.Count, response.Count)
}

func TestFizzBuzzController_GenerateFizzBuzz_ValidationErrors(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIFizzBuzzService(t)
	controller := NewFizzBuzzController(mockService)

	tests := []struct {
		name           string
		request        model.FizzBuzzRequest
		expectedStatus int
	}{
		{
			name: "invalid_int1_zero",
			request: model.FizzBuzzRequest{
				Int1:  0,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_int1_negative",
			request: model.FizzBuzzRequest{
				Int1:  -1,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_int2_zero",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  0,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_int2_negative",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  -1,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_limit_zero",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 0,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_limit_negative",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: -1,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_limit_too_high",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 10001,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_str1_empty",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "",
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_str1_too_long",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  string(make([]byte, 101)), // 101 characters
				Str2:  "buzz",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_str2_empty",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_str2_too_long",
			request: model.FizzBuzzRequest{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  string(make([]byte, 101)), // 101 characters
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request to JSON
			requestBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/fizzbuzz", bytes.NewBuffer(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.GenerateFizzBuzz(c)

			// For validation errors, the handler returns an error, but Echo handles it
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestFizzBuzzController_GenerateFizzBuzz_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIFizzBuzzService(t)
	controller := NewFizzBuzzController(mockService)

	request := model.FizzBuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	mockService.EXPECT().GenerateFizzBuzz(mock.Anything, request).Return(nil, assert.AnError).Once()

	// Marshal request to JSON
	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/fizzbuzz", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GenerateFizzBuzz(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
}

func TestFizzBuzzController_GenerateFizzBuzz_InvalidJSON(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIFizzBuzzService(t)
	controller := NewFizzBuzzController(mockService)

	req := httptest.NewRequest(http.MethodPost, "/fizzbuzz", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GenerateFizzBuzz(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, he.Code)
}

func TestFizzBuzzController_HealthCheck(t *testing.T) {
	e := echo.New()
	mockService := mocks.NewMockIFizzBuzzService(t)
	controller := NewFizzBuzzController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.HealthCheck(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.HealthCheckResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.NotZero(t, response.Timestamp)
}