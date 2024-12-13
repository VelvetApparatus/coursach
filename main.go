package course

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "siaod/course/clock"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/clock"
	station2 "siaod/course/pkg/driverhub"
)

//type TimeTable interface {
//	GetDriveTime(prev, next Point) time.Time
//	GetStationByID(stationID uuid.UUID) Station
//}

// Simulation запускает эмуляцию
func Simulation(ctx context.Context) {
	// Инициализация автобусов
	path := &Path{
		Points: []Point{
			{ID: uuid.New(), Name: "Station A"},
			{ID: uuid.New(), Name: "Station B"},
			{ID: uuid.New(), Name: "Station C"},
		},
		Number: 1,
	}

	bus := &bus.Bus{
		len:    0,
		cap:    50,
		path:   path,
		next:   path.Points[0],
		last:   path.Points[0],
		driver: station2.NewDriverA(),
	}

	// Подписываем автобус на тики часов
	busTick := clock.C().Subscribe()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-busTick:
				if t.After(bus.nextStopTime) {
					err := bus.StopAndDriveNext(nil) // TimeTable пока пустой, но можно настроить
					if err != nil {
						fmt.Printf("Bus error: %v\n", err)
					} else {
						fmt.Printf("Bus moved to: %s at %s\n", bus.last.Name, t.Format("15:04"))
					}
				}
			}
		}
	}()

}
