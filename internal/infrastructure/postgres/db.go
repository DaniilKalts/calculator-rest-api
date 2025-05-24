package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DaniilKalts/calculator-rest-api/internal/domain/models"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Calculation{}); err != nil {
		return nil, err
	}

	return db, nil
}
