package station

import (
	"github.com/google/uuid"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/path"
	"siaod/course/pkg/timetable"
	"sync"
	"time"
)

type BusStation struct {
	station path.Station
	mu      sync.RWMutex
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

func (bst *BusStation) Register(b *bus.Bus) {
	bst.mu.Lock()
	defer bst.mu.Unlock()
	bst.buses[b.ID] = *b
}

func (bst *BusStation) GetNotInWork(
	tt timetable.TimeTable,
	timeTo time.Time,
) *bus.Bus {
	key := bst.getFirst(func(b bus.Bus) bool {
		return !tt.BusOnTheWayToTime(timeTo, b.ID)
	})
	if key == uuid.Nil {
		return nil
	}
	b := bst.GetBus(key)
	return b
}

func (bst *BusStation) getFirst(fn func(b bus.Bus) bool) uuid.UUID {
	bst.mu.RLock()
	defer bst.mu.RUnlock()
	for k, v := range bst.buses {
		if fn(v) {
			return k
		}
	}
	return uuid.Nil
}
