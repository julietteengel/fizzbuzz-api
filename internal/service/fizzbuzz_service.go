package service

import (
	"fmt"
	"strconv"

	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

type fizzBuzzService struct{}

func NewFizzBuzzService() FizzBuzzService {
	return &fizzBuzzService{}
}

func (s *fizzBuzzService) GenerateFizzBuzz(request model.FizzBuzzRequest) (*model.FizzBuzzResponse, error) {
	if request.Limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	if request.Int1 <= 0 || request.Int2 <= 0 {
		return nil, fmt.Errorf("int1 and int2 must be greater than 0")
	}

	if request.Str1 == "" || request.Str2 == "" {
		return nil, fmt.Errorf("str1 and str2 cannot be empty")
	}

	result := make([]string, 0, request.Limit)

	for i := 1; i <= request.Limit; i++ {
		var value string

		isMultipleOfInt1 := i%request.Int1 == 0
		isMultipleOfInt2 := i%request.Int2 == 0

		switch {
		case isMultipleOfInt1 && isMultipleOfInt2:
			value = request.Str1 + request.Str2
		case isMultipleOfInt1:
			value = request.Str1
		case isMultipleOfInt2:
			value = request.Str2
		default:
			value = strconv.Itoa(i)
		}

		result = append(result, value)
	}

	return &model.FizzBuzzResponse{
		Result: result,
		Count:  len(result),
	}, nil
}