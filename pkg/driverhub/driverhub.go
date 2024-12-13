package driverhub

import (
	"github.com/google/uuid"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/timetable"
	"sync"
	"time"
)

type DriverHub struct {
	drivers map[uuid.UUID]driver.Driver
	mu      sync.RWMutex
}

func (dh *DriverHub) GetDriver(id uuid.UUID) driver.Driver {
	dh.mu.RLock()
	defer dh.mu.RUnlock()
	drv, _ := dh.drivers[id]
	return drv
}

func (dh *DriverHub) Register(driver driver.Driver) {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	dh.drivers[driver.ID()] = driver
}

func (dh *DriverHub) GetNotInWork(tt timetable.TimeTable, timeTo time.Time) driver.Driver {
	key := dh.getFirst(func(d driver.Driver) bool {
		return !tt.DriverOnTheWayToTime(timeTo, d.ID())
	})
	if key == uuid.Nil {
		return nil
	}

	drv := dh.drivers[key]
	return drv
}

func (dh *DriverHub) GetFirst(fn func(d driver.Driver) bool) uuid.UUID {
	return dh.getFirst(fn)
}

func (dh *DriverHub) getFirst(fn func(d driver.Driver) bool) uuid.UUID {
	dh.mu.RLock()
	defer dh.mu.RUnlock()
	for k, v := range dh.drivers {
		if fn(v) {
			return k
		}
	}
	return uuid.Nil
}
