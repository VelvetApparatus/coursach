package timetable

import (
	"github.com/google/uuid"
	"siaod/course/pkg/path"
	"time"
)

type TimeTable interface {
	BusOnTheWayToTime(timeTo time.Time, busID uuid.UUID) bool
	DriverOnTheWayToTime(timeTo time.Time, driverID uuid.UUID) bool

	GetBusPositionToTime(busID uuid.UUID, timeTo time.Time) uuid.UUID
	GetDriverPositionToTime(driverID uuid.UUID, timeTo time.Time) uuid.UUID

	GetStationByID(stationID uuid.UUID) path.Station
	GetDriveTime(src path.Point, dest path.Point) time.Duration

	GetPathByID(pathID uuid.UUID) path.Path

	PathsLen() int
	Paths() map[uuid.UUID]path.Path
	GetFirstN(n int, fn func(p path.Path) bool) []uuid.UUID

	GetPathToTime(timeTo time.Time) uuid.UUID

	AssignDriverToPath(pathID uuid.UUID, driverID uuid.UUID)
	AssignBusToPath(pathID uuid.UUID, busID uuid.UUID)
}
