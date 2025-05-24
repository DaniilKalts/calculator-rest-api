package usecase

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"

	"github.com/DaniilKalts/calculator-rest-api/internal/domain/models"
	"github.com/DaniilKalts/calculator-rest-api/internal/infrastructure/postgres/repositories"
)

type CalculationService interface {
	CreateCalculation(expression string) (
		*models.Calculation, error,
	)
	FetchCalculations() ([]models.Calculation, error)
	UpdateCalculation(
		id string, expression string,
	) (*models.Calculation, error)
	DeleteCalculation(id string) error
}

type calculationService struct {
	repo repositories.CalculationRepository
}

func NewCalculationService(repo repositories.CalculationRepository) CalculationService {
	return &calculationService{repo: repo}
}

func (s *calculationService) CreateCalculation(expression string) (
	*models.Calculation, error,
) {
	evaluatedExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}

	result, err := evaluatedExpression.Evaluate(nil)
	if err != nil {
		return nil, err
	}

	newCalculation := &models.Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     fmt.Sprintf("%v", result),
	}

	return s.repo.Create(newCalculation)
}

func (s *calculationService) FetchCalculations() ([]models.Calculation, error) {
	return s.repo.FetchAll()
}

func (s *calculationService) UpdateCalculation(
	id string, expression string,
) (*models.Calculation, error) {
	evaluatedExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}

	result, err := evaluatedExpression.Evaluate(nil)
	if err != nil {
		return nil, err
	}

	updatedCalculation := &models.Calculation{
		ID:         id,
		Expression: expression,
		Result:     fmt.Sprintf("%v", result),
	}

	return s.repo.Update(id, updatedCalculation)
}

func (s *calculationService) DeleteCalculation(id string) error {
	return s.repo.Delete(id)
}
