package stats

import (
	"course/pkg/driver"
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"fmt"
	"math"
	"os"
	"sort"
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
	driversCount        float64
	busCount            float64
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

	s.driversCount = float64(len(dh.Drivers()))

	for _, d := range dh.Drivers() {
		s.averagePathOnDriver += float64(len(tt.GetEach(func(p path.Path) bool { return p.DriverID == d.ID() })))
	}
	s.averagePathOnDriver /= s.driversCount

	adrvs := float64(len(dh.GetEach(func(d driver.Driver) bool { return d.Type() == driver.DriverA })))
	s.drvsDistribution = adrvs / s.driversCount

	s.busCount = float64(len(bs.Buses()))

	for _, b := range bs.Buses() {
		s.averagePathOnBus += float64(len(tt.GetEach(func(p path.Path) bool { return p.BusID == b.ID })))
	}

	s.averagePathOnBus /= float64(s.driversCount)

	ds.exps[optimizer] = append(ds.exps[optimizer], *s)
}

func (ds *DriversStats) SaveStatistics(filename string) error {
	var builder strings.Builder

	// Вывод данных по каждой группе экспериментов
	for optimizer, stats := range ds.exps {
		builder.WriteString(fmt.Sprintf("### Результаты экспериментов для оптимизатора: %s\n\n", optimizer))

		// Подсчет средних значений и дополнительных метрик
		avg := calculateAverage(stats)
		median := calculateMedian(stats)
		variance := calculateVariance(stats)
		stdDev := calculateStdDev(variance)

		// Вывод метрик
		builder.WriteString("#### Сводные метрики:\n\n")
		builder.WriteString(fmt.Sprintf("- Среднее значение:\n"))
		builder.WriteString(fmt.Sprintf("  - Drivers Count: %.2f\n", avg.driversCount))
		builder.WriteString(fmt.Sprintf("  - Bus Count: %.2f\n", avg.busCount))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Driver: %.2f\n", avg.averagePathOnDriver))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Bus: %.2f\n", avg.averagePathOnBus))
		builder.WriteString(fmt.Sprintf("  - Drvs Distribution: %.2f\n\n", avg.drvsDistribution))

		builder.WriteString(fmt.Sprintf("- Медиана:\n"))
		builder.WriteString(fmt.Sprintf("  - Drivers Count: %.2f\n", median.driversCount))
		builder.WriteString(fmt.Sprintf("  - Bus Count: %.2f\n", median.busCount))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Driver: %.2f\n", median.averagePathOnDriver))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Bus: %.2f\n", median.averagePathOnBus))
		builder.WriteString(fmt.Sprintf("  - Drvs Distribution: %.2f\n\n", median.drvsDistribution))

		builder.WriteString(fmt.Sprintf("- Дисперсия:\n"))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Driver: %.2f\n", variance.averagePathOnDriver))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Bus: %.2f\n", variance.averagePathOnBus))
		builder.WriteString(fmt.Sprintf("  - Drvs Distribution: %.2f\n\n", variance.drvsDistribution))

		builder.WriteString(fmt.Sprintf("- Стандартное отклонение:\n"))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Driver: %.2f\n", stdDev.averagePathOnDriver))
		builder.WriteString(fmt.Sprintf("  - Avg Path On Bus: %.2f\n", stdDev.averagePathOnBus))
		builder.WriteString(fmt.Sprintf("  - Drvs Distribution: %.2f\n\n", stdDev.drvsDistribution))

		// Вывод данных по каждому эксперименту
		builder.WriteString("#### Детализация экспериментов:\n\n")
		builder.WriteString(fmt.Sprintf("| Experiment | Drivers Count | Bus Count | Avg Path/Driver | Avg Path/Bus | Drvs Distribution |\n"))
		builder.WriteString(fmt.Sprintf("|------------|---------------|-----------|-----------------|--------------|--------------------|\n"))
		for i, s := range stats {
			builder.WriteString(fmt.Sprintf("| %10d | %.2f | %.2f | %.2f | %.2f | %.2f |\n",
				i+1, s.driversCount, s.busCount, s.averagePathOnDriver, s.averagePathOnBus, s.drvsDistribution))
		}
		builder.WriteString("\n\n")
	}

	// Сохранение в файл
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

// Методы для вычисления статистик
func calculateAverage(stats []stat) stat {
	avg := stat{}
	for _, s := range stats {
		avg.driversCount += s.driversCount
		avg.busCount += s.busCount
		avg.averagePathOnDriver += s.averagePathOnDriver
		avg.averagePathOnBus += s.averagePathOnBus
		avg.drvsDistribution += s.drvsDistribution
	}
	n := float64(len(stats))
	if n > 0 {
		avg.driversCount = avg.driversCount / n
		avg.busCount = avg.busCount / n
		avg.averagePathOnDriver /= n
		avg.averagePathOnBus /= n
		avg.drvsDistribution /= n
	}
	return avg
}

func calculateMedian(stats []stat) stat {
	median := stat{}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].driversCount < stats[j].driversCount
	})
	mid := len(stats) / 2
	if len(stats)%2 == 0 {
		median = stats[mid-1]
	} else {
		median = stats[mid]
	}
	return median
}

func calculateMin(stats []stat) stat {
	min := stats[0]
	for _, s := range stats {
		if s.averagePathOnDriver < min.averagePathOnDriver {
			min = s
		}
	}
	return min
}

func calculateMax(stats []stat) stat {
	max := stats[0]
	for _, s := range stats {
		if s.averagePathOnDriver > max.averagePathOnDriver {
			max = s
		}
	}
	return max
}

func calculateVariance(stats []stat) stat {
	mean := calculateAverage(stats)
	variance := stat{}
	for _, s := range stats {
		variance.averagePathOnDriver += math.Pow(s.averagePathOnDriver-mean.averagePathOnDriver, 2)
	}
	n := float64(len(stats))
	if n > 0 {
		variance.averagePathOnDriver /= n
	}
	return variance
}

func calculateStdDev(variance stat) stat {
	stdDev := stat{}
	stdDev.averagePathOnDriver = math.Sqrt(variance.averagePathOnDriver)
	return stdDev
}
