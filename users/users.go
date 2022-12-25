package users

import (
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
)

type User struct {
	Email string
	Team []drivers.Driver
	Leagues []leagues.League
}