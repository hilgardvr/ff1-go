package drivers

type Driver struct {
	Id     int
	Name   string
	Points int
	Price  int
}

const budget = 1000000

func AssignPrices(drivers []Driver) []Driver {
	var totalPoints int
	var createdDrivers []Driver
	for _, driver := range drivers {
		totalPoints += driver.Points
	}
	for _, driver := range drivers {
		price := driver.Points * budget / totalPoints
		createdDrivers = append(createdDrivers, Driver{
			Id:     driver.Id,
			Name:   driver.Name,
			Points: driver.Points,
			Price:  price,
		})
	}
	return createdDrivers
}
