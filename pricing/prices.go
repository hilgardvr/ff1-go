package pricing

import (
	"hilgardvr/ff1-go/drivers"
	// "hilgardvr/ff1-go/service"
	// "log"
	"math"
)

const budget float64 = 1000000.0
const driversInTeam int = 4
const basePrice float64 = budget * 0.1
const adjustmentFactor float64 = 1.5


func AssignPrices(ds []drivers.Driver) []drivers.Driver {
	var createdDrivers []drivers.Driver
	totalPoints := sumAllDriverPoints(ds)
	for _, driver := range ds {
		price := calcPrice(driver, totalPoints)
		createdDrivers = append(createdDrivers, drivers.Driver{
			Id:     driver.Id,
			Name:   driver.Name,
			Surname: driver.Surname,
			Points: driver.Points,
			Price:  price,
			Constructor: driver.Constructor,
		})
	}
	return createdDrivers
}

func sumAllDriverPoints(drivers []drivers.Driver) int64 {
	var totalPoints int64
	for _, driver := range drivers {
		totalPoints += driver.Points
	}
	return totalPoints
}

func calcPrice(driver drivers.Driver, totalPoints int64) int64 {
	if totalPoints == 0 {
		totalPoints++
	}
	driverPointsShare := float64(driver.Points) / float64(totalPoints)
	price := budget * driverPointsShare * 3 + budget * 0.05
	price = math.Round(price)
	return int64(price)
}