package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/the-great-checkout/transactions-crud/internal/dto"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	FindByID(id uuid.UUID) (*entity.Transaction, error)
	FindAll() ([]entity.Transaction, error)
	Update(transaction *entity.Transaction) error
	Delete(id uuid.UUID) (*entity.Transaction, error)
}

type TransactionService struct {
	repository TransactionRepository
	mapper     TransactionMapper
}

type TransactionMapper interface {
	ToDTO(transaction *entity.Transaction) *dto.Transaction
	FromDTO(transaction *dto.Transaction) *entity.Transaction
}

func NewTransactionService(repository TransactionRepository, mapper TransactionMapper) *TransactionService {
	return &TransactionService{repository: repository, mapper: mapper}
}

func (s *TransactionService) Create(value float64) (*dto.Transaction, error) {
	transaction := &entity.Transaction{
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repository.Create(transaction)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(transaction), nil
}

func (s *TransactionService) GetByID(id uuid.UUID) (*dto.Transaction, error) {
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(transaction), nil
}

func (s *TransactionService) GetAll() ([]dto.Transaction, error) {
	transactions, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.Transaction, len(transactions))
	for i := range transactions {
		dtos[i] = *s.mapper.ToDTO(&transactions[i])
	}

	return dtos, nil
}

func (s *TransactionService) Update(id uuid.UUID, status string, value float64) (*dto.Transaction, error) {
	transaction, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	transaction.Status = entity.Status{Name: status}
	transaction.Value = value
	transaction.UpdatedAt = time.Now()

	err = s.repository.Update(transaction)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(transaction), nil
}

func (s *TransactionService) Delete(id uuid.UUID) (*dto.Transaction, error) {
	transaction, err := s.repository.Delete(id)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(transaction), nil
}
