package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/DaniilKalts/calculator-rest-api/internal/delivery/http"
	"github.com/DaniilKalts/calculator-rest-api/internal/delivery/http/handlers"
	"github.com/DaniilKalts/calculator-rest-api/internal/infrastructure/postgres"
	"github.com/DaniilKalts/calculator-rest-api/internal/infrastructure/postgres/repositories"
	"github.com/DaniilKalts/calculator-rest-api/internal/usecase"
)

func main() {
	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5433 sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

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
