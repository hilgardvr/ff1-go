package users

import (
	"errors"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/races"
)

type User struct {
	Email        string
	TeamName     string
	TeamPriciple string
	Team         []drivers.Driver
	Leagues      []leagues.League
	IsAdmin      bool
	SeasonPoints int
	RacePoints   []races.RacePoints
	Budget       int64
}

func UpdateUserBudgetWithSelections(newDrivers []drivers.Driver, user User) (int64, error) {
	costOldDrivers := 0
	for _, v := range user.Team {
		costOldDrivers += int(v.Price)
	}
	userBudget := user.Budget + int64(costOldDrivers)
	costNewDrivers := 0
	if len(newDrivers) != 4 {
		return user.Budget, errors.New(fmt.Sprint("Invalid number of drivers for user: ", user.Email))
	}
	for _, d := range newDrivers {
		costNewDrivers += int(d.Price)
	}
	if costNewDrivers > int(userBudget) {
		return user.Budget, errors.New(fmt.Sprint("Invalid budget for user: ", user.Email))
	}
	return (userBudget - int64(costNewDrivers)), nil
}
