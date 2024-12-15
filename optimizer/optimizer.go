package optimizer

import (
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable/ttv1"
)

type Optimizer interface {
	Optimize(
		tt *ttv1.TimeTable,
		buses *station.BusStation,
		drvs *driverhub.DriverHub,
	)
}
