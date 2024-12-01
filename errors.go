package course

import "errors"

var (
	BusNotInSafePointErr     = errors.New("bus is not in safe point")
	BusCannotFindNewPointErr = errors.New("bus can not find new point")
	BusReachedEndOfPathErr   = errors.New("bus reaches end of path")
)
