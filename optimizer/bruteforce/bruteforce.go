package bruteforce

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"siaod/course/optimizer"
	bus2 "siaod/course/pkg/bus"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable"
)

type bruteForce struct{}

func NewBrutForceOptimizer() optimizer.Optimizer {
	return &bruteForce{}
}

func (bf *bruteForce) Optimize(
	tt timetable.TimeTable,
	buses *station.BusStation,
	drvs *driverhub.DriverHub,
) {
	for _, p := range tt.Paths() {
		if p.DriverID == uuid.Nil {
			drv := drvs.GetNotInWork(tt, p.StartTime)
			if drv == nil {
				if rand.Int()%2 == 0 {
					drv = driver.NewDriverA()
				} else {
					drv = driver.NewDriverB()
				}
				drvs.Register(drv)
			}

			tt.AssignDriverToPath(p.ID, drv.ID())
		}

		if p.BusID == uuid.Nil {
			bus := buses.GetNotInWork(tt, p.StartTime)
			if bus == nil {
				bus = bus2.NewBus(uuid.New())
				buses.Register(bus)
			}
			tt.AssignBusToPath(p.ID, bus.ID)
		}
	}
}
