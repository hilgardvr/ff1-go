package races

import "hilgardvr/ff1-go/drivers"

type Race struct {
	Race int64
	Season int64
}

type RacePoints struct {
	Race Race
	Drivers []drivers.Driver
	Total int64
}