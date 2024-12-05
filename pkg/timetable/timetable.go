package timetable

import (
	"github.com/google/uuid"
	"siaod/course/pkg/path"
	"sync"
	"time"
)

type TimeTable interface {
	BusOnTheWayToTime(timeTo time.Time, busID uuid.UUID) bool
	GetBusPositionToTime(busID uuid.UUID, timeTo time.Time) uuid.UUID
	GetDriverPositionToTime(driverID uuid.UUID, timeTo time.Time) uuid.UUID
	GetStationByID(stationID uuid.UUID) path.Station
	GetDriveTime(src path.Point, dest path.Point) time.Duration
}

type tt struct {
	mu                sync.RWMutex
	paths             map[uuid.UUID]path.Path
	stationsDistances map[uuid.UUID]map[uuid.UUID]time.Duration
	stations          map[uuid.UUID]path.Station
}

func NewTimeTable() TimeTable { return &tt{} }

func (t *tt) getEachPath(fn func(path path.Path) bool) []uuid.UUID {
	paths := make([]uuid.UUID, 0)

	t.mu.Lock()
	for k, v := range t.paths {
		if fn(v) {
			paths = append(paths, k)
		}
	}
	t.mu.Unlock()
	return paths
}

func (t *tt) getFirstPath(fn func(path path.Path) bool) uuid.UUID {
	t.mu.Lock()
	for k, v := range t.paths {
		if fn(v) {
			return k
		}
	}
	defer t.mu.Unlock()
	return uuid.Nil
}

func (t *tt) compareEach(fn func(a, b path.Path) bool) uuid.UUID {
	t.mu.Lock()
	defer t.mu.Unlock()

	best := uuid.Nil
	for k := range t.paths {
		if best == uuid.Nil {
			best = k
			continue
		}

		if !fn(t.paths[best], t.paths[k]) {
			best = k
		}
	}
	return best
}

func (t *tt) BusOnTheWayToTime(timeTo time.Time, busID uuid.UUID) bool {
	paths := t.getEachPath(func(path path.Path) bool {
		return path.BusID == busID
	})

	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, p := range paths {
		if timeTo.Before(t.paths[p].EndTime) && timeTo.After(t.paths[p].StartTime) {
			return true
		}

	}

	return false
}

func (t *tt) GetBusPositionToTime(busID uuid.UUID, timeTo time.Time) uuid.UUID {
	paths := t.getEachPath(func(path path.Path) bool {
		return path.BusID == busID
	})

	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, pathID := range paths {
		if timeTo.Before(t.paths[pathID].StartTime) || timeTo.After(t.paths[pathID].EndTime) {
			continue
		}
		p := t.paths[pathID]

		tmstp := p.StartTime
		for i := 1; i < len(p.Points); i++ {
			newt := tmstp.Add(t.stationsDistances[p.Points[i-1].ID()][p.Points[i].ID()])
			if timeTo.After(tmstp) && timeTo.Before(newt) {
				return p.Points[i].ID()
			}
			tmstp = newt
		}
	}
	return uuid.Nil
}

func (t *tt) GetDriverPositionToTime(driverID uuid.UUID, timeTo time.Time) uuid.UUID {
	paths := t.getEachPath(func(path path.Path) bool {
		return path.DriverID == driverID
	})

	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, pathID := range paths {
		if timeTo.Before(t.paths[pathID].StartTime) || timeTo.After(t.paths[pathID].EndTime) {
			continue
		}
		p := t.paths[pathID]

		tmstp := p.StartTime
		for i := 1; i < len(p.Points); i++ {
			newt := tmstp.Add(t.stationsDistances[p.Points[i-1].ID()][p.Points[i].ID()])
			if timeTo.After(tmstp) && timeTo.Before(newt) {
				return p.Points[i].ID()
			}
			tmstp = newt
		}
	}
	return uuid.Nil
}

func (t *tt) GetStationByID(stationID uuid.UUID) path.Station {
	t.mu.RLock()
	defer t.mu.RUnlock()
	st := t.stations[stationID]
	return st
}

func (t *tt) GetDriveTime(src path.Point, dest path.Point) time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	dur := t.stationsDistances[src.ID()][dest.ID()]
	return dur
}
