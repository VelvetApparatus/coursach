package main

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"siaod/course/optimizer/bruteforce"
	"siaod/course/optimizer/greedy"
	"siaod/course/pkg/bus"
	"siaod/course/pkg/clock"
	_ "siaod/course/pkg/clock"
	"siaod/course/pkg/driver"
	"siaod/course/pkg/driverhub"
	"siaod/course/pkg/path"
	"siaod/course/pkg/station"
	"siaod/course/pkg/timetable/ttv1"
	"siaod/course/presenter"
	"time"
)

func main() {
	tt, drvs, bss := genScene()

	bfTT := tt.Build()
	bfDrvs := drvs.Build()
	bfBss := bss.Build()
	opt := bruteforce.NewBrutForceOptimizer()
	opt.Optimize(bfTT, bfBss, bfDrvs)
	prs := presenter.Presenter{}
	prs.Present("bruteforce", bfTT, bfDrvs, bfBss)

	gTT := tt.Build()
	gDrvs := drvs.Build()
	gBss := bss.Build()
	opt = greedy.NewGreedyOptimizer()
	opt.Optimize(gTT, gBss, gDrvs)
	prs = presenter.Presenter{}
	prs.Present("greedy", gTT, gDrvs, gBss)

}

func genScene() (*ttv1.TimetableBuilder, *driverhub.DriverHubBuilder, *station.BusStationBuilder) {
	const (
		//driverACount = 5
		driverACount = 1
		//driverBCount = 5
		driverBCount = 1
	)
	bss := []path.Point{
		{Id: uuid.New(), Name: "bus_station_1", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_2", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_3", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_4", IsBusStation: true},
		{Id: uuid.New(), Name: "bus_station_5", IsBusStation: true},
	}

	workStartTime := time.Date(0, 0, 0, 6, 0, 0, 0, time.Local)
	workEndTime := time.Date(0, 0, 0, 23, 0, 0, 0, time.Local)
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

	return tt, hub, stb
}

func genTimeTable(
	busStations []path.Point,
	pathsCount int,
	workStart time.Time,
	workEnd time.Time,
) *ttv1.TimetableBuilder {
	ttb := ttv1.NewBuilder()
	inc := increment()
	rndTime := randTime(workEnd, workStart)
	for i := 0; i < pathsCount; i++ {
		src := rand.IntN(len(busStations))
		dst := rand.IntN(len(busStations))
		stationsCount := rand.IntN(10) + 10

		p := path.NewPath(busStations[src], busStations[dst], inc(), stationsCount, time.Now())
		dstItems := p.GenDstItems()

		var rideDur time.Duration
		for _, item := range dstItems {
			rideDur += item.Dur
		}

		for i := 0; i < 12; i++ {
			p.StartTime = rndTime()
			p.ID = uuid.New()
			p.EndTime = p.StartTime.Add(rideDur)
			ttb.AddPath(p, dstItems)
		}
	}
	return ttb
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
		return time.Date(year, month, day, rand.IntN(fromH-toH)+toH, minArr[rand.IntN(len(minArr))], 0, 0, time.Local)
	}
}
