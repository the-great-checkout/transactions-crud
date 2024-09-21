package controller

import (
	"fmt"
	"net/http"
	_ "time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/the-great-checkout/transactions-crud/internal/dto"
)

type TransactionService interface {
	Create(value float64) (*dto.Transaction, error)
	GetByID(id uuid.UUID) (*dto.Transaction, error)
	GetAll() ([]dto.Transaction, error)
	Update(id uuid.UUID, status string, value float64) (*dto.Transaction, error)
	Delete(id uuid.UUID) (*dto.Transaction, error)
}

type NotificationService interface {
	Publish(message any) error
}

type TransactionController struct {
	transactionService  TransactionService
	notificationService NotificationService
}

func NewTransactionController(
	transactionService TransactionService, notificationService NotificationService) *TransactionController {
	return &TransactionController{
		transactionService:  transactionService,
		notificationService: notificationService,
	}
}

// CreateHandler creates a new transaction
//
//	@Summary		Create a transaction
//	@Description	Create a new transaction with a value
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			transaction	body		dto.Transaction	true	"Transaction Data"
//	@Success		201			{object}	dto.Transaction
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/transactions [post]
func (ctrl *TransactionController) CreateHandler(c echo.Context) error {
	var input dto.Transaction
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	transactionDTO, err := ctrl.transactionService.Create(input.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = ctrl.notificationService.Publish(transactionDTO)
	if err != nil {
		return c.JSON(http.StatusAccepted, transactionDTO)
	}

	return c.JSON(http.StatusCreated, transactionDTO)
}

// GetByIDHandler retrieves a transaction by ID
//
//	@Summary		Get a transaction by ID
//	@Description	Retrieve a single transaction using its ID
//	@Tags			transactions
//	@Produce		json
//	@Param			id	path		string	true	"Transaction ID"
//	@Success		200	{object}	dto.Transaction
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/transactions/{id} [get]
func (ctrl *TransactionController) GetByIDHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	transactionDTO, err := ctrl.transactionService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactionDTO)
}

// GetAllHandler retrieves all transactions
//
//	@Summary		Get all transactions
//	@Description	Retrieve all transactions
//	@Tags			transactions
//	@Produce		json
//	@Success		200	{array}		dto.Transaction
//	@Failure		500	{object}	map[string]string
//	@Router			/transactions [get]
func (ctrl *TransactionController) GetAllHandler(c echo.Context) error {
	transactionsDTOs, err := ctrl.transactionService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactionsDTOs)
}

// UpdateHandler updates a transaction by ID
//
//	@Summary		Update a transaction
//	@Description	Update a transaction's status and value by its ID
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string			true	"Transaction ID"
//	@Param			transaction	body		dto.Transaction	true	"Transaction Data"
//	@Success		200			{object}	dto.Transaction
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/transactions/{id} [put]
func (ctrl *TransactionController) UpdateHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	var input dto.Transaction
	if err = c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updatedTransactionDTO, err := ctrl.transactionService.Update(id, input.Status, input.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = ctrl.notificationService.Publish(updatedTransactionDTO)
	if err != nil {
		fmt.Print(err.Error())
		return c.JSON(http.StatusAccepted, updatedTransactionDTO)
	}

	return c.JSON(http.StatusOK, updatedTransactionDTO)
}

// DeleteHandler deletes a transaction by ID
//
//	@Summary		Delete a transaction
//	@Description	Soft delete a transaction by its ID
//	@Tags			transactions
//	@Produce		json
//	@Param			id	path	string	true	"Transaction ID"
//	@Success		204
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/transactions/{id} [delete]
func (ctrl *TransactionController) DeleteHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	transactionDTO, err := ctrl.transactionService.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = ctrl.notificationService.Publish(transactionDTO)
	if err != nil {
		return c.JSON(http.StatusAccepted, transactionDTO)
	}

	return c.NoContent(http.StatusNoContent)
}
