package mapper

import (
	"github.com/the-great-checkout/transactions-crud/internal/dto"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
)

type TransactionMapper struct {
}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (*TransactionMapper) ToDTO(transaction *entity.Transaction) *dto.Transaction {
	return &dto.Transaction{
		ID:        transaction.ID,
		Status:    transaction.Status.Name,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
		Value:     transaction.Value,
	}
}
func (*TransactionMapper) FromDTO(transaction *dto.Transaction) *entity.Transaction {
	return &entity.Transaction{
		ID: transaction.ID,
		Status: entity.Status{
			Name: transaction.Status,
		},
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
		Value:     transaction.Value,
	}
}
