package driverhub

import (
	"github.com/google/uuid"
	"siaod/course/pkg/driver"
	"sync"
)

type DriverHub struct {
	drivers map[string]driver.Driver
	mu      sync.RWMutex
}

func (dh *DriverHub) GetDriver(id uuid.UUID) driver.Driver {
	dh.mu.RLock()
	defer dh.mu.RUnlock()
	drv, _ := dh.drivers[id.String()]
	return drv
}
