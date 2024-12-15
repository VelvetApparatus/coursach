package presenter

import (
	"course/pkg/driverhub"
	"course/pkg/path"
	"course/pkg/station"
	"course/pkg/timetable/ttv1"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

type Presenter struct {
}

func (p *Presenter) Present(
	filename string,
	tt *ttv1.TimeTable,
	dh *driverhub.DriverHub,
	bst *station.BusStation,
) {
	var builder strings.Builder

	builder.WriteString("# Расписание и сводная информация\n\n")
	builder.WriteString(fmt.Sprintf("## Количество автобусов: %d\n\n", len(bst.Buses())))
	builder.WriteString(fmt.Sprintf("## Количество водителей: %d\n\n", len(dh.Drivers())))
	// Таблица расписания
	builder.WriteString("## Таблица расписания\n\n")
	builder.WriteString("| Путь ID | Номер | Водитель | Автобус | Время начала | Время конца |\n")
	builder.WriteString("|---------|-------|----------|---------|--------------|-------------|\n")
	for id, p := range tt.Paths() {
		builder.WriteString(fmt.Sprintf(
			"| %s | %d | %s | %s | %s | %s |\n",
			id,
			p.Number,
			getDriverName(dh, p.DriverID),
			getBusName(bst, p.BusID),
			p.StartTime.Format("15:04:05"),
			p.EndTime.Format("15:04:05"),
		))
	}

	builder.WriteString("\n")

	// Визуализация рабочего времени водителей
	builder.WriteString("## Рабочее время водителей\n\n")
	builder.WriteString("| Водитель | Рабочие промежутки | Количество путей |\n")
	builder.WriteString("|----------|--------------------|------------------|\n")
	for id, d := range dh.Drivers() {
		paths := tt.GetEach(func(p path.Path) bool { return p.DriverID == id })
		var intervals []string
		for _, pid := range paths {
			p := tt.GetPathByID(pid.ID)
			intervals = append(intervals, fmt.Sprintf("%s-%s", p.StartTime.Format("15:04"), p.EndTime.Format("15:04")))
		}
		builder.WriteString(fmt.Sprintf(
			"| %s | %s | %d |\n",
			d.ID().String(),
			strings.Join(intervals, ", "),
			len(paths),
		))
	}

	builder.WriteString("\n")

	// Количество путей для каждого автобуса
	builder.WriteString("## Количество путей для каждого автобуса\n\n")
	builder.WriteString("| Автобус | Количество путей |\n")
	builder.WriteString("|---------|------------------|\n")
	for id, b := range bst.Buses() {
		paths := tt.GetEach(func(p path.Path) bool { return p.BusID == id })
		builder.WriteString(fmt.Sprintf("| %s | %d |\n", b.ID.String(), len(paths)))
	}

	// Сохранение в файл
	file, err := os.Create(fmt.Sprintf("%s.md", filename))
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(builder.String())
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}

	fmt.Println("Данные успешно сохранены в README.md")

}

// Вспомогательные функции
func getDriverName(dh *driverhub.DriverHub, driverID uuid.UUID) string {
	d := dh.GetDriver(driverID)
	if d == nil {
		return "Не назначен"
	}
	return d.ID().String()
}

func getBusName(bst *station.BusStation, busID uuid.UUID) string {
	b := bst.GetBus(busID)
	if b == nil {
		return "Не назначен"
	}
	return b.ID.String()
}

func main() {
	// Пример использования (инициализация TimeTable, DriverHub и BusStation опущена)
	// Present(timeTable, driverHub, busStation)
}
