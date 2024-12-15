package bruteforce

import (
	"github.com/google/uuid"
	"log/slog"
	"math/rand/v2"
	"siaod/course/optimizer"
	bus2 "siaod/course/pkg/bus"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/path"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable/ttv1"
	"slices"
)

type bruteForce struct{}

func NewBrutForceOptimizer() optimizer.Optimizer {
	return &bruteForce{}
}

func (bf *bruteForce) Optimize(
	tt *ttv1.TimeTable,
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
				slog.Info(
					"brute-force-optimizer",
					slog.String("driverID", drv.ID().String()),
					slog.Int("type", drv.Type()),
					slog.String("status", "hired"),
				)
				drvs.Register(drv)
			}
			paths := tt.GetEach(func(p path.Path) bool {
				return p.DriverID == drv.ID()
			})
			var ps []path.Path
			for _, pp := range paths {
				ps = append(ps, pp)
			}
			slices.SortFunc(ps, func(a, b path.Path) int {
				if a.StartTime.Before(b.StartTime) {
					return -1
				}
				return 1
			})

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
