package driver

import (
	"github.com/google/uuid"
	"siaod/course/pkg/clock"
	"siaod/course/pkg/path"
	"slices"
	"time"
)

type Driver interface {
	ID() uuid.UUID
	NewWorkSession(timeStart, timeEnd time.Time)
	StopWorkSession()
	ActiveToday() bool
	ReadyToWorkNow() bool
	NewDaySession()
	Rest()
	Type() DriverType
	NeedsRest(ps []path.Path) bool
	RestDur() time.Duration
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

func (d *driver) NewWorkSession(timeStart time.Time, timeEnd time.Time) {
	d.timeStart = timeStart
	d.timeEnd = timeEnd
}

func (d *driver) StopWorkSession() {
	d.workCounterH -= int(d.timeEnd.Sub(d.timeStart).Hours())
	if d.workCounterH < 0 {
		d.Rest()
	}
}

func (d *driver) ActiveToday() bool    { return d.active }
func (d *driver) ReadyToWorkNow() bool { return clock.C().Now().After(d.restTimeEnd) }

func (d *driver) NewDaySession() {
	if d.active {
		d.workCounter--
		if d.workCounter < 0 {
			d.active = false
			d.weekendCounter = d.sets.weekendDays
		}
		return
	}
	d.weekendCounter--
	if d.weekendCounter < 0 {
		d.active = true
		d.workCounter = d.sets.workTimeDays
	}
}

func (d *driver) Rest() {
	d.restTimeEnd = clock.C().Now().Add(time.Duration(d.sets.restTimeDur.Milliseconds() / d.sets.restCount))
}

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
