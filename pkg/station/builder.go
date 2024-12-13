package station

import (
	"github.com/google/uuid"
	"maps"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/consts"
	"siaod/course/pkg/timetable"
	"sync"
	"time"
)

type BusStationBuilder struct {
	StationID uuid.UUID
	mu        sync.Mutex
	buses     map[uuid.UUID]bus.Bus

	sets busStationSets
}

type busStationSets struct {
	autoBuy bool
}

func NewBusStationBuilder() *BusStationBuilder {
	return &BusStationBuilder{
		buses: make(map[uuid.UUID]bus.Bus),
		sets: busStationSets{
			autoBuy: true,
		},
	}
}

func (builder *BusStationBuilder) Build(tt timetable.TimeTable) *BusStation {
	station := BusStation{
		station: tt.GetStationByID(builder.StationID),
		buses:   make(map[uuid.UUID]bus.Bus),
	}
	maps.Copy(station.buses, builder.buses)
	return &station
}

func (builder *BusStationBuilder) WithoutAutoBuy() {
	builder.sets.autoBuy = false
}

func (builder *BusStationBuilder) GetFreeBus(
	tt timetable.TimeTable,
	timeTo time.Time,
) (b bus.Bus, err error) {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	for _, v := range builder.buses {
		// проверяем, не имеет ли автобус уже активного рейса
		if tt.BusOnTheWayToTime(timeTo, v.ID) {
			continue
		}
		// если id станций не совпадают -- находится в другом автопарке
		if stationID := tt.GetBusPositionToTime(v.ID, timeTo); stationID != builder.StationID {
			continue
		}
		b = v
		break
	}

	// если по итоги обхода по всем автобусам в таксопарке не нашлось ни одного свободного автобуса
	if b.ID == uuid.Nil {
		// если не можем докупать
		if !builder.sets.autoBuy {
			return bus.Bus{}, consts.NoFreeBussesOnStationError
		}
		busID := builder.buyNewBus()
		b = builder.buses[busID]
	}
	return b, nil
}

func (builder *BusStationBuilder) buyNewBus() uuid.UUID {
	bus := bus.NewBus(uuid.New())
	builder.mu.Lock()
	builder.buses[bus.ID] = *bus
	builder.mu.Unlock()
	return bus.ID
}

// AddBus
// Используем метод в случае если другой автопак купил новый автобус и
// в рамках некоторых рейсов он может оказаться здесь
func (builder *BusStationBuilder) AddBus(bus bus.Bus) {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	builder.buses[bus.ID] = bus
}
