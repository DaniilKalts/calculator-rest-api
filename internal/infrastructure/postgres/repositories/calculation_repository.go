package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/DaniilKalts/calculator-rest-api/internal/domain/models"
)

type CalculationRepository interface {
	Create(calculation *models.Calculation) (*models.Calculation, error)
	FetchAll() ([]models.Calculation, error)
	FetchByID(id string) (*models.Calculation, error)
	Update(id string, calculation *models.Calculation) (
		*models.Calculation, error,
	)
	Delete(id string) error
}

type calculationRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calculationRepository{db: db}
}

func (r *calculationRepository) Create(calculation *models.Calculation) (
	*models.Calculation, error,
) {
	res := r.db.Create(calculation)

	if res.Error != nil {
		return nil, res.Error
	}

	return calculation, nil
}

func (r *calculationRepository) FetchAll() ([]models.Calculation, error) {
	var calculations []models.Calculation
	res := r.db.Find(&calculations)

	if res.Error != nil {
		return nil, res.Error
	}

	return calculations, nil
}

func (r *calculationRepository) FetchByID(id string) (
	*models.Calculation, error,
) {
	var calculation models.Calculation

	err := r.db.First(&calculation, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &calculation, nil
}

func (r *calculationRepository) Update(
	id string,
	calculation *models.Calculation,
) (*models.Calculation, error) {
	existingCalculation, err := r.FetchByID(id)
	if err != nil {
		return nil, err
	}

	existingCalculation.Expression = calculation.Expression
	existingCalculation.Result = calculation.Result

	if err := r.db.Save(existingCalculation).Error; err != nil {
		return nil, err
	}

	return existingCalculation, nil
}

func (r *calculationRepository) Delete(id string) error {
	if err := r.db.Delete(
		&models.Calculation{}, "id = ?", id,
	).Error; err != nil {
		return err
	}

	return nil
}
