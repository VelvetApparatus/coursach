package station

import (
	"course/pkg/bus"
	"github.com/google/uuid"
	"maps"
	"sync"
)

type BusStationBuilder struct {
	mu    sync.Mutex
	buses map[uuid.UUID]bus.Bus

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

func (builder *BusStationBuilder) Build() *BusStation {
	station := BusStation{
		buses: make(map[uuid.UUID]bus.Bus),
	}
	maps.Copy(station.buses, builder.buses)
	return &station
}

func (builder *BusStationBuilder) WithoutAutoBuy() {
	builder.sets.autoBuy = false
}

func (builder *BusStationBuilder) buyNewBus() uuid.UUID {
	bus := bus.NewBus(uuid.New())
	builder.mu.Lock()
	builder.buses[bus.ID] = *bus
	builder.mu.Unlock()
	return bus.ID
}

func (builder *BusStationBuilder) AddBus(bus *bus.Bus) {
	builder.mu.Lock()
	defer builder.mu.Unlock()
	builder.buses[bus.ID] = *bus
}
