package driverhub

import (
	"github.com/google/uuid"
	"maps"
	"siaod/course/pkg/driver"
	"sync"
)

type DriverHubBuilder struct {
	stationID uuid.UUID
	drivers   map[uuid.UUID]driver.Driver
	mu        sync.Mutex

	sets driverHubSets
}

func NewDriverHubBuilder() *DriverHubBuilder {
	return &DriverHubBuilder{
		drivers: make(map[uuid.UUID]driver.Driver),
	}
}

func (dh *DriverHubBuilder) Build() *DriverHub {
	d := DriverHub{
		drivers: make(map[uuid.UUID]driver.Driver),
	}
	maps.Copy(d.drivers, dh.drivers)
	return &d
}

type driverHubSets struct {
	autoHire bool
}

func (dh *DriverHubBuilder) AddDriver(d driver.Driver) {
	dh.mu.Lock()
	defer dh.mu.Unlock()
	dh.drivers[d.ID()] = d
}
