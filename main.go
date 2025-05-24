package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*  Connect to the Database  */

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5433 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

/* ======== Structs ======== */

type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

/* ==== Business Logic ==== */

func calculateExpression(expression string) (string, error) {
	evaluatedExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}

	result, err := evaluatedExpression.Evaluate(nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

/* ==== Handlers ===== */

func getCalculations(c echo.Context) error {
	calculations := []Calculation{}

	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not fetch calculations"},
		)
	}

	return c.JSON(http.StatusOK, calculations)
}

func postCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request payload"},
		)
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid expression"},
		)
	}

	newCalculation := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}
	if err := db.Create(&newCalculation).Error; err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not create calculation"},
		)
	}

	return c.JSON(http.StatusCreated, newCalculation)
}

func patchCalculations(c echo.Context) error {
	id := c.Param("id")

	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request payload"},
		)
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid expression"},
		)
	}

	var calculation Calculation
	if err := db.Find(&calculation, "id = ?", id).Error; err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Could not fetch calculation"},
		)
	}
	calculation.Expression = req.Expression
	calculation.Result = result

	if err := db.Save(&calculation).Error; err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not update calculation"},
		)
	}

	return c.JSON(
		http.StatusOK,
		calculation,
	)
}

func deleteCalculations(c echo.Context) error {
	id := c.Param("id")

	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not delete calculation"},
		)
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	initDB()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8080")
}
