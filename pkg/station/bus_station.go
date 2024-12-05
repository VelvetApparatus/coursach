package station

import (
	"github.com/google/uuid"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/path"
	"sync"
)

type BusStation struct {
	station path.Station
	mu      sync.Mutex
	buses   map[uuid.UUID]bus.Bus
}

func (bst *BusStation) GetBus(id uuid.UUID) *bus.Bus {
	bst.mu.Lock()
	defer bst.mu.Unlock()
	b, ok := bst.buses[id]
	if !ok {
		return nil
	}
	return &b
}
