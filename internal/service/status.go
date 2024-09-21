package service

import (
	"github.com/google/uuid"
	"github.com/the-great-checkout/transactions-crud/internal/dto"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
)

type StatusRepository interface {
	Create(status *entity.Status) error
	FindByID(id uuid.UUID) (*entity.Status, error)
	FindAll() ([]entity.Status, error)
}

type StatusService struct {
	repository StatusRepository
	mapper     StatusMapper
}

type StatusMapper interface {
	ToDTO(status *entity.Status) *dto.Status
	FromDTO(status *dto.Status) *entity.Status
}

func NewStatusService(repository StatusRepository, mapper StatusMapper) *StatusService {
	return &StatusService{repository: repository, mapper: mapper}
}

func (s *StatusService) Create(name string) (*dto.Status, error) {
	status := &entity.Status{
		Name: name,
	}
	err := s.repository.Create(status)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(status), nil
}

func (s *StatusService) GetByID(id uuid.UUID) (*dto.Status, error) {
	status, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(status), nil
}

func (s *StatusService) GetAll() ([]dto.Status, error) {
	statuses, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.Status, len(statuses))
	for i, status := range statuses {
		dtos[i] = *s.mapper.ToDTO(&status)
	}

	return dtos, nil
}
