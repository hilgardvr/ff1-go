package drivers

type Driver struct {
	Id     int
	Name   string
	Points int
	Price  float64
	// RoundBonus  int
}

const budget float64 = 1000000.0
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
			// RoundBonus: 0,
		})
	}
	return createdDrivers
}

func sumAllDriverPoints(drivers []Driver) int {
	var totalPoints int
	for _, driver := range drivers {
		totalPoints += driver.Points
	}
	return totalPoints
}

func calcPrice(driver Driver, totalPoints int) float64 {
	if totalPoints == 0 {
		totalPoints++
	}
	driverPointsShare := float64(driver.Points) / float64(totalPoints)
	price := (driverPointsShare*budget + basePrice) * adjustmentFactor
	return price
}
