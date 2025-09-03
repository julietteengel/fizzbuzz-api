package service

import "github.com/julietteengel/fizzbuzz-api/internal/model"

type FizzBuzzService interface {
	GenerateFizzBuzz(request model.FizzBuzzRequest) (*model.FizzBuzzResponse, error)
}