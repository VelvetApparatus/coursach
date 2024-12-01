package course

import (
	"github.com/google/uuid"
	"siaod/course/clock"
	"sync"
	"time"
)

type PathItem struct {
	TimeBegin time.Time
	TimeEnd   time.Time
	Path      Path
	InWork    bool
}
type TimeTable struct {
	driveTimes map[uuid.UUID]map[uuid.UUID]time.Duration // Время движения между станциями
	stations   map[uuid.UUID]Station                     // Станции по ID
	path       map[uuid.UUID][]PathItem
	pathMu     sync.RWMutex
}

func NewTimeTable() *TimeTable {
	return &TimeTable{
		driveTimes: make(map[uuid.UUID]map[uuid.UUID]time.Duration),
		stations:   make(map[uuid.UUID]Station),
		path:       make(map[uuid.UUID][]PathItem),
	}
}

func (tt *TimeTable) GetNextPathItem() *Path {
	best := new(PathItem)

	tt.pathMu.Lock()
	defer tt.pathMu.Unlock()
	for _, val := range tt.path {
		for i, item := range val {
			if item.InWork {
				continue
			}
			if best == nil || item.TimeBegin.Before(best.TimeBegin) {
				best = &val[i]
			}
		}
	}
	best.InWork = true

	return &best.Path
}

// AddStation добавляет станцию в расписание
func (tt *TimeTable) AddStation(station Station) {
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
func (tt *TimeTable) GetDriveTime(prev, next Point) time.Time {
	duration := tt.driveTimes[prev.ID][next.ID]
	return clock.C().Now().Add(duration)
}

// GetStationByID возвращает станцию по ID
func (tt *TimeTable) GetStationByID(stationID uuid.UUID) Station {
	return tt.stations[stationID]
}
