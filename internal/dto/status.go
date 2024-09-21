package dto

import "github.com/google/uuid"

type Status struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}
