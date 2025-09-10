package service

import (
	"context"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
	"github.com/julietteengel/fizzbuzz-api/internal/repository"
	"strconv"
)

type IFizzBuzzService interface {
	GenerateFizzBuzz(ctx context.Context, request model.FizzBuzzRequest) (*model.FizzBuzzResponse, error)
}

type fizzBuzzService struct {
	statsRepo repository.IStatsRepository
}

func NewFizzBuzzService(statsRepo repository.IStatsRepository) IFizzBuzzService {
	return &fizzBuzzService{
		statsRepo: statsRepo,
	}
}

func (s *fizzBuzzService) GenerateFizzBuzz(ctx context.Context, request model.FizzBuzzRequest) (*model.FizzBuzzResponse, error) {
	result := make([]string, 0, request.Limit)

	for i := 1; i <= request.Limit; i++ {
		var value string
		//
		//Opérateur modulo (%) :
		//- i % request.Int1 = reste de la division de i par request.Int1
		//- Si le reste est 0, alors i est divisible par request.Int1
		// "Est-ce que i divisé par Int1 donne un reste de 0 ?"
		// Si oui → i est un multiple de Int1
		// Si non → i n'est pas un multiple de Int1
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

	// Record request for statistics (async to not block response)
	//- Goroutine pour éviter de bloquer la réponse HTTP avec l'enregistrement des stats
	//- L'enregistrement des stats est un effet de bord non critique
	//- Si la base de données est lente, on ne veut pas ralentir l'API
	// PB: // Sans timeout - goroutine peut rester bloquée indéfiniment
	go func() {
		if err := s.statsRepo.RecordRequest(context.Background(), request); err != nil {
			// Log error but don't fail the request
		}

		// AMELIORATION POSSIBLE :
		//  Avec timeout - goroutine se termine au bout de 5s max
		//go func() {
		//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		//	defer cancel() // Si DB répond jamais → cancel() se déclenche au bout de 5s  => Même si DB down, goroutine meurt après 5s (sinon fuite de ressources)
		//	if err := s.statsRepo.RecordRequest(ctx, request); err != nil {
		//		log.Errorf("Failed to record stats: %v", err)
		//	}
		//}()
	}()

	return &model.FizzBuzzResponse{
		Result: result,
		Count:  len(result),
	}, nil
}
