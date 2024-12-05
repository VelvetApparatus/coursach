package v2

import (
	"github.com/google/uuid"
	"siaod/course"
	"siaod/course/clock"
	"siaod/course/pkg/bus"
	"sync"
	"time"
)

type PathItem struct {
	ID        uuid.UUID
	TimeBegin time.Time
	TimeEnd   time.Time
	Path      course.Path
	DriverID  uuid.UUID
	BusID     uuid.UUID
	InWork    bool
}

type TimeTable struct {
	driveTimes map[uuid.UUID]map[uuid.UUID]time.Duration // Время движения между станциями
	stations   map[uuid.UUID]course.Station              // Станции по id
	path       map[uuid.UUID]PathItem
	pathMu     sync.RWMutex
}

func NewTimeTable() *TimeTable {
	return &TimeTable{
		driveTimes: make(map[uuid.UUID]map[uuid.UUID]time.Duration),
		stations:   make(map[uuid.UUID]course.Station),
		path:       make(map[uuid.UUID][]PathItem),
	}
}

func (tt *TimeTable) GetNextPathItem() *course.Path {
	best := new(PathItem)

	tt.pathMu.Lock()
	defer tt.pathMu.Unlock()
	for _, val := range tt.path {
		if val.InWork {
			continue
		}
		if best == nil || val.TimeBegin.Before(best.TimeBegin) {
			best = &val
		}
	}
	return &best.Path
}

func (tt *TimeTable) SetDriverWithBusOnPath(driver course.Driver, bus bus.Bus, pathID uuid.UUID) {
	tt.pathMu.Lock()
	defer tt.pathMu.Unlock()
	path := tt.path[pathID]
	path.BusID = bus.ID
	path.DriverID = driver.ID()
	path.InWork = true
	tt.path[pathID] = path
}

// AddStation добавляет станцию в расписание
func (tt *TimeTable) AddStation(station course.Station) {
	tt.stations[station.GetID()] = station
}

// AddDriveTime добавляет время движения между станциями
func (tt *TimeTable) AddDriveTime(from, to uuid.UUID, duration time.Duration) {
	if tt.driveTimes[from] == nil {
		tt.driveTimes[from] = make(map[uuid.UUID]time.Duration)
	}
	tt.driveTimes[from][to] = duration
}

// GetDriveTime возвращает время движения между станциями
func (tt *TimeTable) GetDriveTime(prev, next course.Point) time.Time {
	duration := tt.driveTimes[prev.ID][next.ID]
	return clock.C().Now().Add(duration)
}

// GetStationByID возвращает станцию по id
func (tt *TimeTable) GetStationByID(stationID uuid.UUID) course.Station {
	return tt.stations[stationID]
}
