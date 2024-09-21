package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/the-great-checkout/transactions-crud/internal/database"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
	"gorm.io/gorm"
)

type StatusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(postgres database.Postgres) *StatusRepository {
	return &StatusRepository{postgres.DB}
}

func (r *StatusRepository) Create(transaction *entity.Status) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return err
	}

	return nil
}

func (r *StatusRepository) FindByID(id uuid.UUID) (*entity.Status, error) {
	var status entity.Status
	if err := r.db.First(&status, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("status not found")
		}
		return nil, err
	}

	return &status, nil
}

func (r *StatusRepository) FindAll() ([]entity.Status, error) {
	var statuses []entity.Status
	if err := r.db.Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
