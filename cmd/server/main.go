package main

import (
	"log"

	"github.com/DaniilKalts/calculator-rest-api/internal/delivery/http"
	"github.com/DaniilKalts/calculator-rest-api/internal/delivery/http/handlers"
	"github.com/DaniilKalts/calculator-rest-api/internal/infrastructure/postgres"
	"github.com/DaniilKalts/calculator-rest-api/internal/infrastructure/postgres/repositories"
	"github.com/DaniilKalts/calculator-rest-api/internal/usecase"
)

func main() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5433 sslmode=disable"

	db, err := postgres.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	repo := repositories.NewCalculationRepository(db)
	service := usecase.NewCalculationService(repo)
	handler := handlers.NewCalculationHandler(service)

	e := http.NewRouter(handler)

	e.Start("localhost:8080")
}
