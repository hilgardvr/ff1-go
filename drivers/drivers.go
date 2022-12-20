package drivers

import "math"

type Driver struct {
	Id     int     	`json:"id"`
	Name   string  	`json:"name"`
	Points int     	`json:"points"`
	Price  int 		`json:"price"`
}

const budget float64 = 1000000.0
const driversInTeam int = 4
const basePrice float64 = budget * 0.1
const adjustmentFactor float64 = 1.5

func AssignPrices(drivers []Driver) []Driver {
	var createdDrivers []Driver
	totalPoints := sumAllDriverPoints(drivers)
	for _, driver := range drivers {
		price := calcPrice(driver, totalPoints)
		createdDrivers = append(createdDrivers, Driver{
			Id:     driver.Id,
			Name:   driver.Name,
			Points: driver.Points,
			Price:  price,
		})
	}
	return createdDrivers
}

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


func sumAllDriverPoints(drivers []Driver) int {
	var totalPoints int
	for _, driver := range drivers {
		totalPoints += driver.Points
	}
	return totalPoints
}

func calcPrice(driver Driver, totalPoints int) int {
	if totalPoints == 0 {
		totalPoints++
	}
	driverPointsShare := float64(driver.Points) / float64(totalPoints)
	price := (driverPointsShare*budget + basePrice) * adjustmentFactor
	price = math.Round(price)
	return int(price)
}
