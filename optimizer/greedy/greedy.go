package greedy

import (
	"course/optimizer"
	"course/pkg/bus"
	"course/pkg/driver"
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"github.com/google/uuid"
	"log/slog"
)

type greedy struct {
}

func NewGreedyOptimizer() optimizer.Optimizer {
	return &greedy{}
}

func (g *greedy) Optimize(
	tt *ttv1.TimeTable,
	buses *station.BusStation,
	drvs *driverhub.DriverHub,
) {
	observe := int(float64(tt.PathsLen()) * 0.37)
	freePaths := tt.GetFirstN(observe, func(p path.Path) bool { return p.BusID == uuid.Nil && p.DriverID == uuid.Nil })

	for _, pathID := range freePaths {
		p := tt.GetPathByID(pathID)

		b := buses.GetNotInWork(tt, p.StartTime)
		if b == nil {
			b = bus.NewBus(uuid.New())
			buses.Register(b)
		}

		drv := drvs.GetNotInWork(tt, p.StartTime)

		var possibleDrvs []driver.Driver

		switch {
		case drv == nil:
			possibleDrvs = append(possibleDrvs, driver.NewDriverA(), driver.NewDriverB())
			break
		case drv.Type() == driver.DriverA:
			possibleDrvs = append(possibleDrvs, drv)
		case drv.Type() == driver.DriverB:
			possibleDrvs = append(possibleDrvs, drv)
		}
		paths := make([][]path.Path, len(possibleDrvs))
		for i, _ := range possibleDrvs {

			paths[i] = append(paths[i], p)

			item := p
			for {

				// если требуется отдых водителю
				if possibleDrvs[i].NeedsRest(paths[i]) {
					item.EndTime.Add(possibleDrvs[i].RestDur())
				}

				nextPathID := tt.GetPathToTime(item.EndTime)
				if nextPathID == uuid.Nil {
					break
				}

				item = tt.GetPathByID(nextPathID)

				paths[i] = append(paths[i], item)

			}

		}

		bestDriver := possibleDrvs[0]
		bestPathArr := paths[0]
		for i := 1; i < len(paths); i++ {
			if len(paths) > len(bestPathArr) {
				bestPathArr = paths[i]
				bestDriver = possibleDrvs[i]
			}
		}
		if drvs.GetDriver(bestDriver.ID()) == nil {
			slog.Info(
				"greedy-optimizer",
				slog.String("driverID", bestDriver.ID().String()),
				slog.Int("type", bestDriver.Type()),
				slog.String("status", "hired"),
			)
			drvs.Register(bestDriver)
		}

		for _, pdx := range bestPathArr {
			tt.AssignDriverToPath(pdx.ID, bestDriver.ID())
			tt.AssignBusToPath(pdx.ID, b.ID)
		}

	}

	last := tt.GetEach(func(p path.Path) bool {
		return p.BusID == uuid.Nil && p.DriverID == uuid.Nil
	})
	for _, p := range last {
		d := drvs.GetNotInWork(tt, p.StartTime)
		if d == nil {
			d = driver.NewDriverA()
			drvs.Register(d)
		}
		tt.AssignDriverToPath(p.ID, d.ID())

		b := buses.GetNotInWork(tt, p.StartTime)
		if b == nil {
			b = bus.NewBus(uuid.New())
			buses.Register(b)
		}

		tt.AssignBusToPath(p.ID, b.ID)
	}
}
