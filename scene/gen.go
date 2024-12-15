package scene

import (
	"course/config"
	"course/pkg/bus"
	"course/pkg/clock"
	"course/pkg/driver"
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"fmt"
	"github.com/google/uuid"
	"math/rand/v2"
	"time"
)

func GenScene() (*ttv1.TimetableBuilder, *driverhub.DriverHubBuilder, *station.BusStationBuilder) {
	bss := make([]path.Point, 0, config.C().InitialBusStationsCount)
	for i := 0; i < config.C().InitialBusStationsCount; i++ {
		bss = append(bss, path.Point{
			Id:           uuid.New(),
			Name:         fmt.Sprintf("BusStation%d", i+1),
			IsBusStation: true,
		})
	}

	workStartTime := time.Date(0, 0, 0, 6, 0, 0, 0, time.Local)
	workEndTime := time.Date(0, 0, 0, 23, 0, 0, 0, time.Local)

	tt := genTimeTable(bss, config.C().DistinctPathCount, workStartTime, workEndTime)
	stb := station.NewBusStationBuilder()
	for range bss {
		for i := 0; i < config.C().InitialBusCount; i++ {
			stb.AddBus(bus.NewBus(uuid.New()))
		}
	}

	hub := driverhub.NewDriverHubBuilder()

	for i := 0; i < config.C().InitialDriverACount; i++ {
		hub.AddDriver(driver.NewDriverA())
	}
	for i := 0; i < config.C().InitialDriverBCount; i++ {
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

		for i := 0; i < config.C().TimeSeriesPathsCount; i++ {
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
