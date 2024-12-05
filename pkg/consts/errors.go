package consts

import "errors"

var (
	NoFreeBussesOnStationError  = errors.New("no free bus on station")
	NoFreeDriversOnStationError = errors.New("no free driverhub on station")
	BusNotInSafePointErr        = errors.New("bus is not in safe point")
	BusCannotFindNewPointErr    = errors.New("bus can not find new point")
	BusReachedEndOfPathErr      = errors.New("bus reaches end of path")
)
