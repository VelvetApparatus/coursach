package optimizer

import (
	"course/pkg/driverhub"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
)

type Optimizer interface {
	Optimize(
		tt *ttv1.TimeTable,
		buses *station.BusStation,
		drvs *driverhub.DriverHub,
	)
}
