package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/the-great-checkout/transactions-crud/internal/dto"
)

type StatusService interface {
	Create(name string) (*dto.Status, error)
	GetByID(id uuid.UUID) (*dto.Status, error)
	GetAll() ([]dto.Status, error)
}

type StatusController struct {
	statusService StatusService
}

func NewStatusController(
	statusService StatusService) *StatusController {
	return &StatusController{
		statusService: statusService,
	}
}

func (ctrl *StatusController) CreateHandler(c echo.Context) error {
	var input dto.Status
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	status, err := ctrl.statusService.Create(input.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, status)
}

func (ctrl *StatusController) GetByIDHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	status, err := ctrl.statusService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, status)
}

func (ctrl *StatusController) GetAllHandler(c echo.Context) error {
	statuses, err := ctrl.statusService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, statuses)
}
