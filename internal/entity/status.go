package entity

import "github.com/google/uuid"

type Status struct {
	ID   uuid.UUID `bson:"_id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name string    `bson:"name" gorm:"type:varchar(255);not null;unique"`
}
