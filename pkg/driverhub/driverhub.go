package driverhub

import (
	"course/pkg/driver"
	"course/pkg/path"
	"course/pkg/timetable/ttv1"
	"github.com/google/uuid"
	"maps"
	"sync"
	"time"
)

type DriverHub struct {
	drivers map[uuid.UUID]driver.Driver
	mu      sync.RWMutex
}

func (dh *DriverHub) Drivers() map[uuid.UUID]driver.Driver {
	res := make(map[uuid.UUID]driver.Driver)
	dh.mu.RLock()
	maps.Copy(res, dh.drivers)
	defer dh.mu.RUnlock()
	return res
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

func (dh *DriverHub) GetNotInWork(
	tt *ttv1.TimeTable,
	timeTo time.Time,
) driver.Driver {
	drvs := dh.GetEach(func(d driver.Driver) bool {
		return !tt.DriverOnTheWayToTime(timeTo, d.ID())
	})

	if len(drvs) == 0 {
		return nil
	}

	for _, drv := range drvs {
		ps := tt.GetEach(func(p path.Path) bool {
			return p.DriverID == drv.ID()
		})

		var pss []path.Path
		for _, v := range ps {
			if v.EndTime.Before(timeTo) {
				pss = append(pss, v)
			}
		}

		if drv.NeedsRest(pss) {
			return nil
		}

		return drv
	}
	return nil
}

func (dh *DriverHub) GetEach(fn func(d driver.Driver) bool) map[uuid.UUID]driver.Driver {
	res := make(map[uuid.UUID]driver.Driver)
	dh.mu.RLock()
	defer dh.mu.RUnlock()
	for _, v := range dh.drivers {
		if fn(v) {
			res[v.ID()] = v
		}
	}
	return res
}

func (dh *DriverHub) GetFirst(fn func(d driver.Driver) bool) uuid.UUID {
	dh.mu.RLock()
	defer dh.mu.RUnlock()
	for k, v := range dh.drivers {
		if fn(v) {
			return k
		}
	}
	return uuid.Nil
}
