package course

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"siaod/course/clock"
	_ "siaod/course/clock"
	"time"
)

type Point struct {
	ID   uuid.UUID
	Name string
}

type Station interface {
	GetID() uuid.UUID
	To() int
	From() int
}
type Path struct {
	Points  []Point
	Number  int
	PathDur time.Duration
}

func (p *Point) GetID() uuid.UUID {
	return p.ID
}
func (p *Path) Last() Point {
	return p.Points[len(p.Points)-1]
}

func (p *Path) GetNext(point Point) Point {
	for i := len(p.Points) - 1; i > 0; i-- {
		if p.Points[i].ID == point.ID {
			return p.Points[i-1]
		}
	}
	return Point{}
}

type Driver interface {
	ID() uuid.UUID
	NewWorkSession(timeStart, timeEnd time.Time)
	StopWorkSession()
	ActiveToday() bool
	ReadyToWorkNow() bool
	NewDaySession()
	Rest()
}

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

	bus := &Bus{
		len:    0,
		cap:    50,
		path:   path,
		next:   path.Points[0],
		last:   path.Points[0],
		driver: NewDriverA(),
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
