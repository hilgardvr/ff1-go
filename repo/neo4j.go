package repo

import (
	"encoding/csv"
	"errors"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"
	"os"
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.Driver
}

var driverData = []drivers.Driver{}

func (n *Neo4jRepo)Init(config *config.Config) error {
	driver, err := neo4j.NewDriver(
		config.Neo4jUri,
		neo4j.BasicAuth(
			config.Neo4jUsername,
			config.Neo4jPassword,
			"",
		),
	)
	if err != nil {
		return err
	}
	err = driver.VerifyConnectivity()

	if err != nil {
		return err
	}
	n.driver = driver
	driverData, err = createDriverData()
	return err
}

func readDriverData() ([][]string, error) {
	path, err := getDriverFilePath()
	if err != nil {
		return [][]string{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	return data, err
}

func createDriverData() ([]drivers.Driver, error) {
	driverData, err := readDriverData()
	if err != nil {
		return []drivers.Driver{}, err
	}
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
	
func (n Neo4jRepo)GetDrivers() []drivers.Driver {
	return driverData
}

func (n Neo4jRepo)AddUser(user users.User) (users.User, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			merge (u:User {email: $email})
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	email, found := result.(map[string]interface{})["email"] 
	if !found {
		return users.User{}, err
	}
	u := users.User{
		Email: email.(string),
	}
	return u, err
}

func (n Neo4jRepo) SetLoginCode(email string, generatedCode string) (string, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			merge (u:User {email: $email})
			set u.loginCode = $generatedCode
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"generatedCode": generatedCode,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	code, found := result.(map[string]interface{})["loginCode"] 
	if !found {
		return "", err
	}
	c := code.(string)
	return c, err
}

func (n Neo4jRepo) DeleteLoginCode(email string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			set u.loginCode = null
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	return err
}

func (n Neo4jRepo) ValidateLoginCode(email string, codeToTest string) bool {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email, loginCode: $loginCode})
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"loginCode": codeToTest,
		})
		record, err := result.Single()
		return record, err
	})
	return err == nil
}

func (n Neo4jRepo) SaveSession(email, uuid string, duration time.Duration) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			merge (s:Session {uuid: $uuid})
			merge (u)-[h:HAS_SESSION]->(s)
			set s.createdAt = timestamp()
			set s.expiry = timestamp() + $expiry
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"uuid": uuid,
			"expiry": duration.Milliseconds(),
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	return err
}

func (n Neo4jRepo) GetSession(uuid string) (string, bool) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (s:Session {uuid: $uuid})
			where s.expiry > timestamp()
			match (u:User)-[:HAS_SESSION]->(s)
			return u { .* } as user
		`,
		map[string]interface{}{
			"uuid": uuid,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	if err != nil {
		return "", false
	}
	res, found := result.(map[string]interface{})["email"] 
	if !found {
		return "", false
	}
	email := res.(string)
	n.DeleteLoginCode(email)
	return email, found
}

func (n Neo4jRepo) SaveTeam(user users.User, selectedDrivers []drivers.Driver) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	var drivers []int
	for _, d := range selectedDrivers {
		drivers = append(drivers, d.Id)
	}
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			merge (t:Team {drivers: $team})
			merge (u)-[h:HAS_TEAM]->(t)
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": user.Email,
			"team": drivers,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
	})
	return err
}

func (n Neo4jRepo) GetTeam(user users.User) ([]drivers.Driver, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})-[:HAS_TEAM]->(t:Team)
			return t.drivers as team
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		record, err := result.Single()
		if err != nil {
			return []drivers.Driver{}, err
		}
		team, found := record.Get("team")
		if !found {
			return []drivers.Driver{}, errors.New("Could not find user in result")
		}
		return team, nil
	})
	if err != nil {
		return []drivers.Driver{}, nil
	}
	res, found := result.([]interface{})
	if !found {
		return []drivers.Driver{}, errors.New("Could not find the team in db results")
	}
	var drivers []drivers.Driver
	for _, driverId := range res {
		for _, d := range driverData {
			if d.Id == int(driverId.(int64)) {
				drivers = append(drivers, d)
			}
		}
	}
	return drivers, err
}

func (n Neo4jRepo) DeleteTeam(user users.User) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})-[:HAS_TEAM]->(t:Team)
			detach delete t
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		resultSummary, err := result.Consume()
		if err != nil {
			return []drivers.Driver{}, err
		}
		return resultSummary, err
	})
	return err
}