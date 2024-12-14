package ttv1

import (
	"github.com/google/uuid"
	"maps"
	"siaod/course/pkg/path"
	"siaod/course/pkg/timetable"
	"sync"
	"time"
)

type Config struct {
	InitBusCount      int
	InitDriversACount int
	InitDriversBCount int
}

type TimetableBuilder struct {
	paths             map[uuid.UUID]path.Path
	stationsDistances map[uuid.UUID]map[uuid.UUID]time.Duration
	stations          map[uuid.UUID]path.Station
	mu                sync.Mutex
}

func NewBuilder() *TimetableBuilder {
	return &TimetableBuilder{}
}

func (builder *TimetableBuilder) Build() timetable.TimeTable {
	t := tt{}
	maps.Copy(builder.stations, t.stations)
	maps.Copy(builder.paths, t.paths)
	maps.Copy(builder.stationsDistances, t.stationsDistances)
	return &t
}

func (builder *TimetableBuilder) AddPath(p path.Path, dstItems []path.DstItem) {
	if builder.pathExists(p.ID) {
		return
	}

	for _, point := range p.Points {
		if builder.stationExists(point.ID()) {
			continue
		}
		builder.addStation(&point)
	}

	for _, dstItem := range dstItems {
		builder.addDistance(dstItem.From, dstItem.To, dstItem.Dur)
	}

	for i := 1; i < len(p.Points); i++ {
		if !builder.dstExists(p.Points[i-1].ID(), p.Points[i].ID()) {
			builder.addDistance(p.Points[i-1].ID(), p.Points[i].ID(), time.Minute*5)
		}
	}

}

func (builder *TimetableBuilder) stationExists(stationID uuid.UUID) bool {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	_, ok := builder.stations[stationID]
	return ok
}

func (builder *TimetableBuilder) pathExists(pathID uuid.UUID) bool {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	_, ok := builder.paths[pathID]
	return ok
}

func (builder *TimetableBuilder) dstExists(src uuid.UUID, dst uuid.UUID) bool {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	_, way1 := builder.stationsDistances[src][dst]
	_, way2 := builder.stationsDistances[dst][src]
	return way1 && way2
}

func (builder *TimetableBuilder) addStation(station path.Station) {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	builder.stations[station.ID()] = station
	builder.stationsDistances[station.ID()] = make(map[uuid.UUID]time.Duration)
}

func (builder *TimetableBuilder) addDistance(
	src uuid.UUID,
	dst uuid.UUID,
	distance time.Duration,
) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	builder.stationsDistances[src][dst] = distance
	builder.stationsDistances[dst][src] = distance
}

func (builder *TimetableBuilder) addPath(
	path path.Path,
) {
	builder.mu.Lock()
	builder.paths[path.ID] = path
	builder.mu.Unlock()
}
