package course

import "math/rand/v2"

type station struct {
	Point
	PathN    int
	Velocity int
}

func (s *station) From() int {
	return rand.IntN(s.PathN) + s.Velocity
}

func (s *station) To() int {
	return rand.IntN(s.PathN) + s.Velocity
}

func NewStation(point Point) Station {
	return &station{
		Point:    point,
		PathN:    0,
		Velocity: 0,
	}
}
