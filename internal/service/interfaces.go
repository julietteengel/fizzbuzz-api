package service

import (
	"context"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

type IFizzBuzzService interface {
	GenerateFizzBuzz(ctx context.Context, request model.FizzBuzzRequest) (*model.FizzBuzzResponse, error)
}