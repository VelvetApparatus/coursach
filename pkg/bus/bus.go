package bus

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"siaod/course/pkg/clock"
	"siaod/course/pkg/consts"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/path"
	"siaod/course/pkg/timetable"
	"time"
)

type Bus struct {
	ID  uuid.UUID
	len int
	cap int

	onRide       bool
	driver       driver.Driver
	path         *path.Path
	next         path.Point
	nextStopTime time.Time

	last         path.Point
	lastStopTime time.Time
}

func NewBus(id uuid.UUID) *Bus {
	return &Bus{
		ID: id,
	}
}

func (b *Bus) GetPosition() path.Point {
	return b.last
}

func (b *Bus) StopAndDriveNext(tt timetable.TimeTable) error {

	if b.next.ID() == b.path.Last().ID() {
		return consts.BusReachedEndOfPathErr
	}

	nextStopTime := b.GetStopTime(tt)
	b.last = b.next
	b.lastStopTime = clock.C().Now()

	b.ServeStation(tt.GetStationByID(b.last.ID()))

	b.next = b.path.GetNext(b.last)
	if b.next.ID() == uuid.Nil {
		return consts.BusCannotFindNewPointErr
	}

	b.nextStopTime = nextStopTime

	return nil
}

func (b *Bus) ServeStation(station path.Station) {
	b.toStation(station)
	b.fromStation(station)
}

func (b *Bus) fromStation(station path.Station) {
	b.len = max(b.cap, b.len+station.From())
}

func (b *Bus) toStation(station path.Station) {
	b.len = min(b.len-station.To(), 0)
}

func (b *Bus) StartPathDrive(ctx context.Context, tt timetable.TimeTable) error {
	b.onRide = true
	ticker := clock.C().Subscribe()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case timestamp := <-ticker:
			if timestamp.After(b.nextStopTime) {
				err := b.StopAndDriveNext(tt)
				if err != nil {
					if errors.Is(err, consts.BusReachedEndOfPathErr) {
						b.onRide = false
						return nil
					}
					return err
				}
			}
		}
	}
}

func (b *Bus) GetStopTime(tt timetable.TimeTable) time.Time {
	return clock.C().Now().Add(tt.GetDriveTime(b.last, b.next))
}

func (b *Bus) ChangePath(path *path.Path) error {
	if b.path != nil {
		if b.last.ID() != b.path.Last().ID() {
			return consts.BusNotInSafePointErr
		}
	}
	b.path = path
	return nil
}

func (b *Bus) SwapDriver(newDriver driver.Driver) {
	b.driver = newDriver
}
