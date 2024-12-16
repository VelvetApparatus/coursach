package ttv1

import (
	"course/pkg/path"
	"github.com/google/uuid"
	"sync"
	"time"
)

type TimeTable struct {
	mu                sync.RWMutex
	paths             map[uuid.UUID]path.Path
	stationsDistances map[uuid.UUID]map[uuid.UUID]time.Duration
	stations          map[uuid.UUID]path.Station
}

func (t *TimeTable) Paths() map[uuid.UUID]path.Path {
	return t.paths
}
func (t *TimeTable) PathsLen() int {
	return len(t.paths)
}

func (t *TimeTable) GetFirstN(n int, fn func(p path.Path) bool) []uuid.UUID {
	res := make([]uuid.UUID, 0, n)
	t.mu.RLock()
	defer t.mu.RUnlock()
	for k, v := range t.paths {
		if fn(v) {
			res = append(res, k)
			if len(res) == n {
				break
			}
		}
	}
	return res
}

func (t *TimeTable) GetPathToTime(timeTo time.Time) uuid.UUID {
	ks := t.getEachPath(func(path path.Path) bool {
		return path.BusID == uuid.Nil &&
			path.DriverID == uuid.Nil &&
			path.StartTime.After(timeTo) &&
			path.StartTime.Before(timeTo.Add(time.Minute*30))
	})
	if len(ks) == 0 {
		return uuid.Nil
	}

	bestKey := ks[0]

	t.mu.RLock()
	for i := 1; i < len(ks); i++ {
		if t.paths[ks[i]].StartTime.Sub(timeTo) >
			t.paths[bestKey].StartTime.Sub(timeTo) {
			bestKey = ks[i]
		}
	}
	t.mu.RUnlock()

	return bestKey
}

func (t *TimeTable) getEachPath(fn func(path path.Path) bool) []uuid.UUID {
	paths := make([]uuid.UUID, 0)

	t.mu.RLock()
	defer t.mu.RUnlock()
	for k, v := range t.paths {
		if fn(v) {
			paths = append(paths, k)
		}
	}
	return paths
}

func (t *TimeTable) BusOnTheWayToTime(timeTo time.Time, busID uuid.UUID) bool {
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

func (t *TimeTable) DriverOnTheWayToTime(timeTo time.Time, driverID uuid.UUID) bool {
	paths := t.getEachPath(func(path path.Path) bool {
		return path.BusID == driverID
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

func (t *TimeTable) AssignDriverToPath(pathID uuid.UUID, driverID uuid.UUID) {
	t.mu.Lock()
	defer t.mu.Unlock()
	p := t.paths[pathID]
	p.DriverID = driverID
	t.paths[pathID] = p
}

func (t *TimeTable) AssignBusToPath(pathID uuid.UUID, busID uuid.UUID) {
	t.mu.Lock()
	defer t.mu.Unlock()
	p := t.paths[pathID]
	p.BusID = busID
	t.paths[pathID] = p
}

func (t *TimeTable) GetPathByID(pathID uuid.UUID) path.Path {
	t.mu.RLock()
	defer t.mu.RUnlock()
	p := t.paths[pathID]
	return p
}

func (t *TimeTable) GetEach(fn func(p path.Path) bool) map[uuid.UUID]path.Path {
	res := make(map[uuid.UUID]path.Path)
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, p := range t.paths {
		if fn(p) {
			res[p.ID] = p
		}
	}
	return res
}
