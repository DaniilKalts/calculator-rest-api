package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CalculationExpression struct {
	ID         int    `json:"id"`
	Expression string `json:"expression"`
	Result     int    `json:"result"`
}

var calculations = []CalculationExpression{}

func getCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)

	e.Start("localhost:8080")
}
