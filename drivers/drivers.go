package drivers

import (
	"hilgardvr/ff1-go/constructor"
	// "math"
)

type Driver struct {
	Id     int64     	`json:"id"`
	Name   string  	`json:"name"`
	Surname string  	`json:"surname"`
	Points int64     	`json:"points"`
	Price  int64 		`json:"price"`
	Constructor constructor.Constructor `json:"constructor"`
}


func ValidateTeam(drivers []Driver, userBudget int64) bool {
	sum := 0
	if len(drivers) != 4 {
		return false
	}
	for _, d := range drivers {
		sum += int(d.Price)
	}
	if sum > int(userBudget) {
		return false
	}
	return true
}
