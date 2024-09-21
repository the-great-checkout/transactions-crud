package mapper

import (
	"github.com/the-great-checkout/transactions-crud/internal/dto"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
)

type StatusMapper struct {
}

func NewStatusMapper() *StatusMapper {
	return &StatusMapper{}
}

func (*StatusMapper) ToDTO(transaction *entity.Status) *dto.Status {
	return &dto.Status{
		ID:   transaction.ID,
		Name: transaction.Name,
	}
}
func (*StatusMapper) FromDTO(transaction *dto.Status) *entity.Status {
	return &entity.Status{
		ID:   transaction.ID,
		Name: transaction.Name,
	}
}
