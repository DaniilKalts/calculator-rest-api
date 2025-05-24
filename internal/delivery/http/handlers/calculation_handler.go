package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DaniilKalts/calculator-rest-api/internal/domain/models"
	"github.com/DaniilKalts/calculator-rest-api/internal/usecase"
)

type CalculationHandler struct {
	service usecase.CalculationService
}

func NewCalculationHandler(service usecase.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: service}
}

func (h *CalculationHandler) HandleCreateCalculation(ctx echo.Context) error {
	var req models.CalculationRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request payload"},
		)
	}

	calculation, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not create calculation"},
		)
	}

	return ctx.JSON(http.StatusCreated, calculation)
}

func (h *CalculationHandler) HandleFetchCalculation(ctx echo.Context) error {
	calculations, err := h.service.FetchCalculations()

	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not fetch calculations"},
		)
	}

	return ctx.JSON(http.StatusOK, calculations)
}

func (h *CalculationHandler) HandleUpdateCalculation(ctx echo.Context) error {
	id := ctx.Param("id")

	var req models.CalculationRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request payload"},
		)
	}

	calculation, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not create calculation"},
		)
	}

	return ctx.JSON(http.StatusOK, calculation)
}

func (h *CalculationHandler) HandleDeleteCalculation(ctx echo.Context) error {
	id := ctx.Param("id")

	err := h.service.DeleteCalculation(id)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Could not delete calculation"},
		)
	}

	return ctx.NoContent(http.StatusNoContent)
}
