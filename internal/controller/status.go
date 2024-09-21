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

// CreateHandler creates a new status
//
//	@Summary		Create a status
//	@Description	Create a new status with a name
//	@Tags			statuses
//	@Accept			json
//	@Produce		json
//	@Param			status	body		dto.Status	true	"Status Data"
//	@Success		201		{object}	dto.Status
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/statuses [post]
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

// GetByIDHandler retrieves a status by ID
//
//	@Summary		Get a status by ID
//	@Description	Retrieve a single status using its ID
//	@Tags			statuses
//	@Produce		json
//	@Param			id	path		string	true	"Status ID"
//	@Success		200	{object}	dto.Status
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/statuses/{id} [get]
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

// GetAllHandler retrieves all statuses
//
//	@Summary		Get all statuses
//	@Description	Retrieve all statuses
//	@Tags			statuses
//	@Produce		json
//	@Success		200	{array}		dto.Status
//	@Failure		500	{object}	map[string]string
//	@Router			/statuses [get]
func (ctrl *StatusController) GetAllHandler(c echo.Context) error {
	statuses, err := ctrl.statusService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, statuses)
}
