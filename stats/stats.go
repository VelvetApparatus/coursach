package stats

import (
	"course/pkg/driver"
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"fmt"
	"os"
	"strings"
)

type DriversStats struct {
	exps map[string][]stat
}

func NewDriversStats() *DriversStats {
	return &DriversStats{
		exps: make(map[string][]stat),
	}
}

type stat struct {
	driversCount        int
	busCount            int
	averagePathOnDriver float64
	averagePathOnBus    float64
	// отношение DriverA к всем остальным
	drvsDistribution float64
}

func (ds *DriversStats) Collect(
	tt *ttv1.TimeTable,
	dh *driverhub.DriverHub,
	bs *station.BusStation,
	optimizer string,
) {
	if ds.exps[optimizer] == nil {
		ds.exps[optimizer] = make([]stat, 0)
	}
	s := new(stat)

	s.driversCount = len(dh.Drivers())

	for _, d := range dh.Drivers() {
		s.averagePathOnDriver += float64(len(tt.GetEach(func(p path.Path) bool { return p.DriverID == d.ID() })))
	}
	s.averagePathOnDriver /= float64(s.driversCount)

	adrvs := float64(len(dh.GetEach(func(d driver.Driver) bool { return d.Type() == driver.DriverA })))
	s.drvsDistribution = float64(s.driversCount) / adrvs

	s.busCount = len(bs.Buses())

	for _, b := range bs.Buses() {
		s.averagePathOnBus += float64(len(tt.GetEach(func(p path.Path) bool { return p.BusID == b.ID })))
	}

	s.averagePathOnBus /= float64(s.driversCount)

	ds.exps[optimizer] = append(ds.exps[optimizer], *s)
}

func (ds *DriversStats) SaveStatistics(filename string) error {

	var builder strings.Builder

	// Подсчет средних значений
	averageStats := stat{}
	experimentCount := 0

	for _, stats := range ds.exps {
		for _, s := range stats {
			averageStats.driversCount += s.driversCount
			averageStats.busCount += s.busCount
			averageStats.averagePathOnDriver += s.averagePathOnDriver
			averageStats.averagePathOnBus += s.averagePathOnBus
			averageStats.drvsDistribution += s.drvsDistribution
			experimentCount++
		}
	}

	if experimentCount > 0 {
		averageStats.driversCount /= experimentCount
		averageStats.busCount /= experimentCount
		averageStats.averagePathOnDriver /= float64(experimentCount)
		averageStats.averagePathOnBus /= float64(experimentCount)
		averageStats.drvsDistribution /= float64(experimentCount)
	}

	// Вывод средних значений
	builder.WriteString(fmt.Sprintf("## Средние значения по всем экспериментам:\n\n"))
	builder.WriteString(fmt.Sprintf("Drivers Count: %d\n\n", averageStats.driversCount))
	builder.WriteString(fmt.Sprintf("Bus Count: %d\n\n", averageStats.busCount))
	builder.WriteString(fmt.Sprintf("Average Path On Driver: %.2f\n\n", averageStats.averagePathOnDriver))
	builder.WriteString(fmt.Sprintf("Average Path On Bus: %.2f\n\n", averageStats.averagePathOnBus))
	builder.WriteString(fmt.Sprintf("Drivers Distribution: %.2f\n\n", averageStats.drvsDistribution))

	// Вывод данных по каждому эксперименту
	for optimizer, stats := range ds.exps {
		builder.WriteString(fmt.Sprintf("#### Результаты экспериментов для оптимизатора: %s\n\n", optimizer))
		builder.WriteString(fmt.Sprintf("Experiment | Drivers Count | Bus Count | Avg Path/Driver | Avg Path/Bus | Drvs Distribution\n"))
		builder.WriteString(fmt.Sprintf("-----------|---------------|-----------|-----------------|--------------|------------------|\n"))
		for i, s := range stats {
			builder.WriteString(fmt.Sprintf("| %10d | %14d | %9d | %15.2f | %13.2f | %17.2f | \n",
				i+1, s.driversCount, s.busCount, s.averagePathOnDriver, s.averagePathOnBus, s.drvsDistribution))
		}
		builder.WriteString("\n\n")
	}
	file, err := os.Create(fmt.Sprintf("%s.md", filename))
	if err != nil {
		return fmt.Errorf("Ошибка создания файла: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(builder.String())
	if err != nil {
		return fmt.Errorf("Ошибка записи в файл: %w", err)
	}

	return nil
}
