package path

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Point struct {
	Id           uuid.UUID
	Name         string
	IsBusStation bool
}

func (p *Point) ID() uuid.UUID {
	return p.Id
}

func (p *Point) IsEnd() bool {
	return p.IsBusStation
}

type Station interface {
	ID() uuid.UUID
	IsEnd() bool
}
type Path struct {
	ID      uuid.UUID
	Points  []Point
	Number  int
	PathDur time.Duration

	BusID     uuid.UUID
	DriverID  uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

func (p *Path) IsPlanned() bool { return p.BusID != uuid.Nil && p.DriverID != uuid.Nil }

func (p *Path) Last() Station {
	return &p.Points[len(p.Points)-1]
}

func (p *Path) GenDstItems() []DstItem {
	dstis := make([]DstItem, len(p.Points)-1)
	for i := 1; i < len(p.Points)-1; i++ {
		dstis[i] = DstItem{
			To:   p.Points[i].ID(),
			From: p.Points[i-1].ID(),
			Dur:  time.Duration(rand.Intn(12)+3) * time.Minute,
		}
	}
	return dstis
}

func NewPath(
	src, dst Point,
	number int,
	stations int,
	startTime time.Time,
) Path {
	points := make([]Point, stations)
	for i := 0; i < stations; i++ {
		points[i] = Point{
			Id:   uuid.New(),
			Name: randomName(12),
		}
	}
	points[0] = src
	points[len(points)-1] = dst

	return Path{
		ID:        uuid.New(),
		Number:    number,
		StartTime: startTime,
		Points:    points,
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomName(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type DstItem struct {
	To   uuid.UUID
	From uuid.UUID
	Dur  time.Duration
}
