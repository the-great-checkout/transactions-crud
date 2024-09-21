package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID        uuid.UUID      `bson:"_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
	DeletedAt gorm.DeletedAt `bson:"deleted_at" gorm:"index"`
	IsDeleted bool           `bson:"is_deleted" gorm:"default:false"`
	StatusID  uuid.UUID      `bson:"-" gorm:"type:uuid"`
	Status    Status         `bson:"status" gorm:"foreignKey:StatusID;references:ID"`
	Value     float64        `bson:"value" gorm:"default:0;notnull"`
}
