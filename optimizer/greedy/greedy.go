package greedy

import (
	"github.com/google/uuid"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/path"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable"
)

type greedy struct {
}

func (g *greedy) Optimize(
	tt timetable.TimeTable,
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

		driverAID := drvs.GetFirst(func(d driver.Driver) bool {
			return d.ReadyToWorkNow() && d.Type() == driver.DriverA
		})
		driverBID := drvs.GetFirst(func(d driver.Driver) bool {
			return d.ReadyToWorkNow() && d.Type() == driver.DriverB
		})

		var possibleDrvs []driver.Driver

		switch {
		case driverAID == uuid.Nil && driverBID != uuid.Nil:
			possibleDrvs = append(possibleDrvs, drvs.GetDriver(driverBID))
		case driverAID != uuid.Nil && driverBID == uuid.Nil:
			possibleDrvs = append(possibleDrvs, drvs.GetDriver(driverAID))
		case driverBID == uuid.Nil && driverAID == uuid.Nil:
			possibleDrvs = append(possibleDrvs, driver.NewDriverA(), driver.NewDriverB())
		default:
			possibleDrvs = append(possibleDrvs, drvs.GetDriver(driverAID), drvs.GetDriver(driverBID))
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

		drvs.Register(bestDriver)

		for _, pdx := range bestPathArr {
			tt.AssignDriverToPath(pdx.ID, bestDriver.ID())
			tt.AssignBusToPath(pdx.ID, b.ID)
		}

	}
}
