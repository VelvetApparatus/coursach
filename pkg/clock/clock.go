package clock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var c *Clock

func C() *Clock { return c }

func init() {
	c = NewClock(time.Date(2024, time.November, 30, 6, 0, 0, 0, time.UTC), time.Second*4)
	go c.Start(context.Background())
}

// Clock - глобальные часы для синхронизации всех событий
type Clock struct {
	mu          sync.Mutex
	currentTime time.Time
	tick        time.Duration
	subscribers []chan time.Time
}

func NewClock(start time.Time, tick time.Duration) *Clock {
	return &Clock{
		currentTime: start,
		tick:        tick,
	}
}

// Subscribe позволяет подписаться на события тика часов
func (c *Clock) Subscribe() chan time.Time {
	ch := make(chan time.Time)
	c.subscribers = append(c.subscribers, ch)
	return ch
}

// Start запускает часы
func (c *Clock) Start(ctx context.Context) {
	ticker := time.NewTicker(c.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Clock stopped")
			for _, ch := range c.subscribers {
				close(ch)
			}
			return
		case <-ticker.C:
			c.mu.Lock()
			c.currentTime = c.currentTime.Add(c.tick)
			c.mu.Unlock()
			fmt.Printf("Clock tick: %s\n", c.currentTime.Format("15:04"))
			for _, ch := range c.subscribers {
				ch <- c.currentTime
			}
		}
	}
}

func (c *Clock) Now() time.Time {
	for !c.mu.TryLock() {
		res := c.currentTime
		c.mu.Unlock()
		return res
	}
	return time.Time{}
}
