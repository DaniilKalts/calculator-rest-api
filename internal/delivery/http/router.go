package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/DaniilKalts/calculator-rest-api/internal/delivery/http/handlers"
)

func NewRouter(handler *handlers.CalculationHandler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", handler.HandleFetchCalculation)
	e.POST("/calculations", handler.HandleCreateCalculation)
	e.PATCH("/calculations/:id", handler.HandleUpdateCalculation)
	e.DELETE("/calculations/:id", handler.HandleDeleteCalculation)

	return e
}
