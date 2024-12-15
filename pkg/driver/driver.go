package driver

import (
	"github.com/google/uuid"
	"siaod/course/pkg/path"
	"slices"
	"time"
)

type Driver interface {
	ID() uuid.UUID
	Type() DriverType
	NeedsRest(ps []path.Path) bool
	RestDur() time.Duration
	WorkDur() time.Duration
}

type DriverType = int

const (
	DriverA = iota
	DriverB
)

type driver struct {
	id           uuid.UUID
	timeStart    time.Time
	timeEnd      time.Time
	restTimeEnd  time.Time
	workCounterH int

	workCounter    int
	weekendCounter int
	active         bool

	sets driverSets
}

type driverSets struct {
	restTimeDur  time.Duration
	workTimeDur  time.Duration
	restCount    int64
	workTimeDays int
	weekendDays  int
	typ          DriverType
}

func NewDriverA() Driver {
	return newDriver(driverSets{
		restTimeDur:  time.Hour,
		workTimeDur:  time.Hour * 8,
		restCount:    1,
		workTimeDays: 5,
		weekendDays:  2,
		typ:          DriverA,
	})
}

func NewDriverB() Driver {
	return newDriver(driverSets{
		restTimeDur:  4 * time.Hour,
		workTimeDur:  24 * time.Hour,
		restCount:    12,
		workTimeDays: 5,
		weekendDays:  2,
		typ:          DriverB,
	})
}

func (d *driver) ID() uuid.UUID { return d.id }

func newDriver(sets driverSets) Driver { return &driver{id: uuid.New(), sets: sets} }

func (d *driver) Type() DriverType { return d.sets.typ }

func (d *driver) NeedsRest(ps []path.Path) bool {
	slices.SortFunc(ps, func(a, b path.Path) int {
		if a.StartTime.Before(b.StartTime) {
			return 1
		}
		return -1
	})
	var timeInWork time.Duration
	for _, p := range ps {
		// time in drive
		timeInWork += p.EndTime.Sub(p.StartTime)
		if timeInWork > d.sets.workTimeDur {
			return true
		}
	}
	return false
}

func (d *driver) RestDur() time.Duration {
	return d.sets.restTimeDur
}

func (d *driver) WorkDur() time.Duration { return d.sets.workTimeDur }
