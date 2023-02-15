package users

import (
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/races"
)

type User struct {
	Email string
	TeamName string
	TeamPriciple string
	Team []drivers.Driver
	Leagues []leagues.League
	IsAdmin bool
	SeasonPoints int
	RacePoints []races.RacePoints
	Budget int64
	Picks int64
}