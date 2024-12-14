package main

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"siaod/course/optimizer/bruteforce"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/clock"
	_ "siaod/course/pkg/clock"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/path"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable"
	"siaod/course/pkg/timetable/ttv1"
	"time"
)

func main() {
	tt, drvs, bss := genScene()
	opt := bruteforce.NewBrutForceOptimizer()
	opt.Optimize(tt, bss, drvs)
}

func genScene() (timetable.TimeTable, *driverhub.DriverHub, *station.BusStation) {
	const (
		driverACount = 15
		driverBCount = 15
	)
	bss := []path.Point{
		{Id: uuid.New(), Name: "bus_station_1", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_2", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_3", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_4", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_5", IsBusStation: true},
	}

	workStartTime := time.Date(0, 0, 0, 6, 0, 0, 0, time.Local)
	workEndTime := time.Date(0, 0, 0, 24, 0, 0, 0, time.Local)
	tt := genTimeTable(bss, 15, workStartTime, workEndTime)
	stb := station.NewBusStationBuilder()
	for range bss {
		for i := 0; i < 5; i++ {
			stb.AddBus(bus.NewBus(uuid.New()))
		}
	}

	hub := driverhub.NewDriverHubBuilder()

	for i := 0; i < driverACount; i++ {
		hub.AddDriver(driver.NewDriverA())
	}
	for i := 0; i < driverBCount; i++ {
		hub.AddDriver(driver.NewDriverB())
	}

	return tt, hub.Build(), stb.Build()
}

func genTimeTable(
	busStations []path.Point,
	pathsCount int,
	workStart time.Time,
	workEnd time.Time,
) timetable.TimeTable {
	ttb := ttv1.NewBuilder()
	inc := increment()
	rndTime := randTime(workStart, workEnd)
	for i := 0; i < pathsCount; i++ {
		src := rand.IntN(len(busStations))
		dst := rand.IntN(len(busStations))
		stationsCount := rand.IntN(15 + 1)

		p := path.NewPath(busStations[src], busStations[dst], inc(), stationsCount, time.Now())
		dstItems := p.GenDstItems()
		var rideDur time.Duration
		for _, item := range dstItems {
			rideDur += item.Dur
		}

		for i := 0; i < 7; i++ {
			p.StartTime = rndTime()
			p.ID = uuid.New()
			p.EndTime = p.StartTime.Add(rideDur)
			ttb.AddPath(p, dstItems)
		}
	}
	return ttb.Build()
}

func increment() func() int {
	inc := 0

	return func() int {
		inc++
		return inc
	}
}

func randTime(from time.Time, to time.Time) func() time.Time {
	fromH := from.Hour()
	toH := to.Hour()
	minArr := []int{5, 10, 15, 20, 25, 30, 35, 45, 50, 55, 0}
	return func() time.Time {
		year, month, day := clock.C().Now().Date()
		return time.Date(year, month, day, rand.IntN(toH-fromH)+fromH, minArr[rand.IntN(len(minArr))], 0, 0, time.Local)
	}
}
