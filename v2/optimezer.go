package v2

import (
	"github.com/google/uuid"
	"siaod/course"
	"siaod/course/pkg/bus"
)

type Optimizer struct {
}

func Optimize(
	paths []PathItem,
	drivers map[uuid.UUID]course.Driver,
	buses map[uuid.UUID]bus.Bus,
) (*TimeTable, error) {
	tt := new(TimeTable)
	observeEdge := int(float64(len(drivers)) * 0.37)
	for i := 0; i < observeEdge; i++ {
		path := paths[i]
		tt.SetDriverWithBusOnPath(drivers[i], buses[i], path.ID)
	}
}
