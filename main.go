package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

var calculations = []Calculation{}

func getCalculations(c echo.Context) error {
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
	calculations = append(calculations, newCalculation)

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

	for i := 0; i < len(calculations); i++ {
		if calculations[i].ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
}

func deleteCalculations(c echo.Context) error {
	id := c.Param("id")

	for idx, calculation := range calculations {
		if calculation.ID == id {
			calculations = append(calculations[:idx], calculations[idx+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8080")
}
