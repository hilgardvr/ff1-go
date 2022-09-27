package repo

import (
	"encoding/csv"
	"errors"
	"hilgardvr/ff1-go/drivers"
	"os"
	"strconv"
)

var parsedDrivers []drivers.Driver

const driverFile = "/repo/data/drivers.csv"

func Init() error {
	path, err := getDriverFilePath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	driverData, err := createDrivers(data)
	if err != nil {
		return err
	}
	parsedDrivers = driverData
	return nil
}

func GetDrivers() []drivers.Driver {
	dst := make([]drivers.Driver, len(parsedDrivers))
	copy(dst, parsedDrivers)
	return dst
}

func createDrivers(driverData [][]string) ([]drivers.Driver, error) {
	// return driverData
	var allDrivers []drivers.Driver
	for _, line := range driverData {
		if len(line) != 3 {
			return allDrivers, errors.New("Driver data unexpected format")
		}
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return allDrivers, err
		}
		name := line[1]
		points, err := strconv.Atoi(line[2])
		if err != nil {
			return allDrivers, err
		}
		driver := drivers.Driver{
			Id:     id,
			Name:   name,
			Points: points,
			Price:  0,
		}
		allDrivers = append(allDrivers, driver)
	}
	allDrivers = drivers.AssignPrices(allDrivers)
	return allDrivers, nil
}

func getDriverFilePath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	joined := path + driverFile

	return joined, nil
}
