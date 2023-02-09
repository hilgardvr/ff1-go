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

const budget float64 = 1000000.0
const driversInTeam int = 4
const basePrice float64 = budget * 0.1
const adjustmentFactor float64 = 1.5

// func AssignPrices(drivers []Driver) []Driver {
// 	// currentRaces, err := service.GetAllRacesForCurrentSeason()
// 	// if err != nil {
// 	// 	return []Driver{}
// 	// }
// 	var createdDrivers []Driver
// 	totalPoints := sumAllDriverPoints(drivers)
// 	for _, driver := range drivers {
// 		price := calcPrice(driver, totalPoints, len(currentRaces))
// 		createdDrivers = append(createdDrivers, Driver{
// 			Id:     driver.Id,
// 			Name:   driver.Name,
// 			Surname: driver.Surname,
// 			Points: driver.Points,
// 			Price:  price,
// 			Constructor: driver.Constructor,
// 		})
// 	}
// 	return createdDrivers
// }

func ValidateTeam(drivers []Driver) bool {
	sum := 0
	if len(drivers) != 4 {
		return false
	}
	for _, d := range drivers {
		sum += int(d.Price)
	}
	if sum > int(budget) {
		return false
	}
	return true
}


// func sumAllDriverPoints(drivers []Driver) int64 {
// 	var totalPoints int64
// 	for _, driver := range drivers {
// 		totalPoints += driver.Points
// 	}
// 	return totalPoints
// }

// func calcPrice(driver Driver, totalPoints int64, numberOfRaces int) int64 {
// 	if totalPoints == 0 {
// 		totalPoints++
// 	}
// 	// driverPointsShare := float64(driver.Points) / float64(totalPoints)
// 	// price := (driverPointsShare*budget + basePrice) * adjustmentFactor
// 	driverPointsShare := float64(driver.Points) / float64(totalPoints)
// 	price := budget * driverPointsShare
// 	price = math.Round(price)
// 	return int64(price)
// }
