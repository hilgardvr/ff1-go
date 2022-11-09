package repo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"
	"math/rand"
	"os"
	"strconv"
)

const driverFile = "/repo/data/drivers.csv"

type LocalFileSystemRepo struct {
 	parsedDrivers []drivers.Driver
	users []users.User
	//uuid - email
 	sessions map[string]string
	loginCodes map[string]string
}

func (l *LocalFileSystemRepo) Init(config *config.Config) error {
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
	l.parsedDrivers = driverData
	l.sessions = map[string]string{}
	l.loginCodes = map[string]string{}
	return nil
}

func (l LocalFileSystemRepo) GetDrivers() []drivers.Driver {
	dst := make([]drivers.Driver, len(l.parsedDrivers))
	copy(dst, l.parsedDrivers)
	return dst
}

func (l LocalFileSystemRepo) AddUser(u users.User) (users.User, error) {
	l.users = append(l.users, u)
	return u, nil
}


func (l LocalFileSystemRepo) SetLoginCode(email string) string {
	if code, found := l.loginCodes[email]; found {
		return code
	}
	code := rand.Intn(100000)
	str := strconv.Itoa(code)
	padded := fmt.Sprintf("%05s", str)
	l.loginCodes[email] = padded
	fmt.Println("logincodes:", l.loginCodes)
	return padded
}

func (l LocalFileSystemRepo) DeleteLoginCode(email string) {
	delete(l.loginCodes, email)
}

func (l LocalFileSystemRepo) ValidateLoginCode(email string, codeToTest string) bool {
	if code, found := l.loginCodes[email]; found {
		return code == codeToTest
	}
	return false
}

func (l LocalFileSystemRepo) SaveSession(email, uuid string) error {
	l.sessions[uuid] = email
	return nil
}

func (l LocalFileSystemRepo) GetSession(uuid string) (string, bool) {
	email := l.sessions[uuid]
	if email == "" {
		return "", false
	}
	return email, true
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
