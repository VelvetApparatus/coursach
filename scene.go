package course

import (
	"context"
	"github.com/google/uuid"
	"log"
	"siaod/course/clock"
)

type Scene struct {
	buses   map[uuid.UUID]Bus
	drivers map[uuid.UUID]Driver
	tt      *TimeTable
}

func (s *Scene) Start(ctx context.Context) error {
	ticker := clock.C().Subscribe()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case timestamp := <-ticker:
			for busID, bus := range s.buses {
				if bus.onRide {
					continue
				}

				path := s.tt.GetNextPathItem()

				for k, drv := range s.drivers {
					if !drv.ActiveToday() {
						continue
					}
					if !drv.ReadyToWorkNow() {
						continue
					}

					s.drivers[k].NewWorkSession(timestamp, timestamp.Add(path.PathDur))
					bus.SwapDriver(drv)
					bus.onRide = true
					break
				}

				if !bus.onRide {
					log.Printf("cannot find new driver")
					continue
				}

				if path != nil {
					*s.buses[busID].path = *path
					go func() {
						err := bus.StartPathDrive(ctx, s.tt)
						if err != nil {
							log.Fatal(err.Error())
							return
						}

					}()
				}
			}
		}
	}
}
