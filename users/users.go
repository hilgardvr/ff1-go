package users

import (
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/races"
)

type User struct {
	Email string
	Team []drivers.Driver
	Leagues []leagues.League
	IsAdmin bool
	SeasonPoints int
	RacePoints []races.RacePoints
}