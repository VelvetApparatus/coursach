package optimizer

import (
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable"
)

type Optimizer interface {
	Optimize(
		tt timetable.TimeTable,
		buses *station.BusStation,
		drvs *driverhub.DriverHub,
	)
}
