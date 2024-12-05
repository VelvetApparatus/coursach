package path

import (
	"github.com/google/uuid"
	"time"
)

type Point struct {
	id           uuid.UUID
	Name         string
	IsBusStation bool
}

func (p *Point) ID() uuid.UUID {
	return p.id
}

func (p *Point) To() int {
	return 0
}

func (p *Point) From() int {
	return 0
}

func (p *Point) IsEnd() bool {
	return p.IsBusStation
}

type Station interface {
	ID() uuid.UUID
	To() int
	From() int
	IsEnd() bool
}
type Path struct {
	ID           uuid.UUID
	Points       []Point
	Number       int
	PathDur      time.Duration
	RideInterval time.Duration

	BusID     uuid.UUID
	DriverID  uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

func (p *Path) IsPlanned() bool { return p.BusID != uuid.Nil && p.DriverID != uuid.Nil }

func (p *Path) Last() Station {
	return &p.Points[len(p.Points)-1]
}

func (p *Path) GetNext(point Point) Point {
	for i := len(p.Points) - 1; i > 0; i-- {
		if p.Points[i].ID() == point.ID() {
			return p.Points[i-1]
		}
	}
	return Point{}
}
