package driverhub

import (
	"github.com/google/uuid"
	"siaod/course/pkg/consts"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/timetable"
	"sync"
	"time"
)

type DriverHubBuilder struct {
	stationID uuid.UUID
	drivers   map[uuid.UUID]driver.Driver
	mu        sync.Mutex

	sets driverHubSets
}

type driverHubSets struct {
	autoHire bool
}

func (dh *DriverHubBuilder) GetFreeDriver(
	tt timetable.TimeTable,
	timeTo time.Time,
) (drv driver.Driver, err error) {
	dh.mu.Lock()
	for _, d := range dh.drivers {
		if !d.ActiveToday() {
			continue
		}
		if !d.ReadyToWorkNow() {
			continue
		}

		if tt.GetDriverPositionToTime(d.ID(), timeTo) != dh.stationID {
			continue
		}
		drv = d

		break
	}
	dh.mu.Unlock()

	if drv == nil {
		if !dh.sets.autoHire {
			return nil, consts.NoFreeDriversOnStationError
		}
		driverID := dh.hireNewDriver()
		drv = dh.drivers[driverID]
	}

	return drv, nil
}

func (dh *DriverHubBuilder) AddDriver(d driver.Driver) {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	dh.drivers[d.ID()] = d
}

func (dh *DriverHubBuilder) hireNewDriver() uuid.UUID {
	drv := driver.NewDriverA()
	dh.mu.Lock()
	defer dh.mu.Unlock()
	dh.drivers[drv.ID()] = drv
	return drv.ID()
}
