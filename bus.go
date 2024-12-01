package course

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"siaod/course/clock"
	"time"
)

type Bus struct {
	ID  uuid.UUID
	len int
	cap int

	onRide       bool
	driver       Driver
	path         *Path
	next         Point
	nextStopTime time.Time

	last         Point
	lastStopTime time.Time
}

func (b *Bus) StopAndDriveNext(tt *TimeTable) error {

	if b.next.ID == b.path.Last().ID {
		return BusReachedEndOfPathErr
	}

	nextStopTime := b.GetStopTime(tt)
	b.last = b.next
	b.lastStopTime = clock.C().Now()

	b.ServeStation(tt.GetStationByID(b.last.ID))

	b.next = b.path.GetNext(b.last)
	if b.next.ID == uuid.Nil {
		return BusCannotFindNewPointErr
	}

	b.nextStopTime = nextStopTime

	return nil
}

func (b *Bus) ServeStation(station Station) {
	b.toStation(station)
	b.fromStation(station)
}

func (b *Bus) fromStation(station Station) {
	b.len = max(b.cap, b.len+station.From())
}

func (b *Bus) toStation(station Station) {
	b.len = min(b.len-station.To(), 0)
}

func (b *Bus) StartPathDrive(ctx context.Context, tt *TimeTable) error {
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
					if errors.Is(err, BusReachedEndOfPathErr) {
						b.onRide = false
						return nil
					}
					return err
				}
			}
		}
	}
}

func (b *Bus) GetStopTime(tt *TimeTable) time.Time {
	return tt.GetDriveTime(b.last, b.next)
}

func (b *Bus) ChangePath(path *Path) error {
	if b.path != nil {
		if b.last.ID != b.path.Last().ID {
			return BusNotInSafePointErr
		}
	}
	b.path = path
	return nil
}

func (b *Bus) SwapDriver(newDriver Driver) {
	b.driver = newDriver
}
