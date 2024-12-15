package bus

import (
	"github.com/google/uuid"
)

type Bus struct {
	ID uuid.UUID
}

func NewBus(id uuid.UUID) *Bus {
	return &Bus{
		ID: id,
	}
}
