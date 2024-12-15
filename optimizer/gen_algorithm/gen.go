package gen_algorithm

import (
	"course/optimizer"
	"course/pkg/bus"
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"github.com/google/uuid"
	"slices"
	"time"
)

type ga struct {
	opt optimizer.Optimizer
}

func New(opt optimizer.Optimizer) optimizer.Optimizer {
	return &ga{opt: opt}
}

func (g *ga) Optimize(
	tt *ttv1.TimeTable,
	buses *station.BusStation,
	drvs *driverhub.DriverHub,
) {
	// используем оптимизатор, чтобы построить первичные пути
	g.opt.Optimize(tt, buses, drvs)

	// 100 эпох
	epoch := epochCounter(100)
	_, next := epoch()
	for next {
		_, next = epoch()
		g.do(tt, buses, drvs)
	}
}

func (g *ga) do(
	tt *ttv1.TimeTable,
	buses *station.BusStation,
	drvs *driverhub.DriverHub,
) {
	parents := g.selectForCrossover(buses, drvs, tt, 10)
	index := 0
	for index+1 < len(parents) {
		g.crossover(tt, parents[index], parents[index+1])
		index += 2
	}
}

func (g *ga) crossover(
	tt *ttv1.TimeTable,
	p1, p2 bus.Bus,
) {
	p1Paths := tt.GetEach(func(p path.Path) bool { return p.BusID == p1.ID })
	p2Paths := tt.GetEach(func(p path.Path) bool { return p.BusID == p2.ID })

	genom := len(p1Paths) + len(p2Paths)

	counter := min(genom, len(p1Paths))
	for _, v := range p1Paths {
		tt.AssignBusToPath(v.ID, p1.ID)
		counter--
		if counter == 0 {
			break
		}
	}

	counter = min(genom, len(p2Paths))

	for _, v := range p2Paths {
		tt.AssignBusToPath(v.ID, p2.ID)
		counter--
		if counter == 0 {
			break
		}
	}

}

func (g *ga) calcFitness(
	tt *ttv1.TimeTable,
	dh *driverhub.DriverHub,
	bus bus.Bus,
) int {
	drvs := make(map[uuid.UUID]struct{})
	penalty := 0

	ps := tt.GetEach(func(p path.Path) bool { return p.BusID == bus.ID })

	for _, p := range ps {
		if _, ok := drvs[p.DriverID]; ok {
			d := dh.GetDriver(p.DriverID)
			drvPaths := tt.GetEach(func(p path.Path) bool { return p.DriverID == d.ID() })

			var (
				workTimeAll       time.Duration
				workTimeAfterRest time.Duration
			)

			for _, drvp := range drvPaths {
				workTimeAll += drvp.EndTime.Sub(drvp.StartTime)
				if drvp.EndTime.Before(p.StartTime) &&
					drvp.StartTime.After(p.StartTime.Add(-d.WorkDur())) {
					workTimeAfterRest += drvp.EndTime.Sub(drvp.StartTime)
				}
			}
			if workTimeAll > d.WorkDur() {
				penalty += 500
			}
			if workTimeAfterRest > d.WorkDur() {
				penalty += 250
			}

			drvs[p.DriverID] = struct{}{}
		}
	}

	return len(drvs)*1000 - len(ps)*500 + penalty
}

func (g *ga) selectForCrossover(
	bs *station.BusStation,
	dh *driverhub.DriverHub,
	tt *ttv1.TimeTable,
	n int,
) []bus.Bus {
	bsMap := bs.Buses()
	sel := make([]bus.Bus, 0, len(bsMap))
	n = min(n, len(bsMap))

	for _, v := range bsMap {
		sel = append(sel, v)
	}
	slices.SortFunc(sel, func(a, b bus.Bus) int {
		if g.calcFitness(tt, dh, a) < g.calcFitness(tt, dh, b) {
			return -1
		}
		return 1
	})

	return sel[:n]
}

func epochCounter(n int) func() (int, bool) {
	ep := n

	return func() (int, bool) {
		ep--
		return ep, ep > 0
	}
}
