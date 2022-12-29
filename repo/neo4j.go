package repo

import (
	"encoding/csv"
	"errors"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/users"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.Driver
}

var driverData = []drivers.Driver{}

func migrate() error {
	fs, err := ioutil.ReadDir("./repo/data")
	if err != nil {
		log.Println("Error opening directory:", err)
		return err
	}
	log.Println(fs)
	return nil
}

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
	err = migrate()
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
	
func (n Neo4jRepo)GetDrivers() ([]drivers.Driver, error) {
	return driverData, nil
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
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

func (n Neo4jRepo) SaveLeague(user users.User, leagueName string, passcode string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	league, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			merge (l:League {name: $leagueName, passcode: $passcode})
			merge (u)-[:LEAGUE]->(l)
			return l { .* } as league
		`,
		map[string]interface{}{
			"email": user.Email,
			"leagueName": leagueName,
			"passcode": passcode,
		})
		record, err := result.Single()
		if err != nil {
			return leagues.League{}, err
		}
		league, found := record.Get("league")
		if !found {
			return leagues.League{}, errors.New("Could not find league in result")
		}
		return league, nil
	})
	log.Println("League created:", league)
	return err
}

func (n Neo4jRepo) GetLeagueForUser(user users.User) ([]leagues.League, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})-[:LEAGUE]-(l:League) 
			return l { .* } as league
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		if err != nil {
			return []leagues.League{}, err
		}
		var ls []leagues.League
		for result.Next() {
			record := result.Record()
			league, found := record.Get("league")
			if found {
				r := league.(map[string]interface{})
				l := leagues.League{
					Name: r["name"].(string),
					Passcode: r["passcode"].(string),
				}
				ls = append(ls, l)
			}
		}
		return ls, err
	})
	ls, found := result.([]leagues.League)
	if !found {
		return []leagues.League{}, errors.New("Could not find the team in db results")
	}
	return ls, err
}

func (n Neo4jRepo) JoinLeague(user users.User, passcode string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email}) 
			match (l:League {passcode: $passcode})
			merge (u)-[:LEAGUE]->(l)
			return l { .* } as league
		`,
		map[string]interface{}{
			"email": user.Email,
			"passcode": passcode,
		})
		if err != nil {
			return []leagues.League{}, err
		}
		record, err := result.Single()
		if err != nil {
			return leagues.League{}, err
		}
		league, found := record.Get("league")
		if !found {
			return leagues.League{}, errors.New("Could not find league in result")
		}
		return league, err
	})
	return err
}

func (n Neo4jRepo) GetLeagueMembers(leaguePasscode string) ([]users.User, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (l:League {passcode: $passcode})-[:LEAGUE]-(u:User)
			return u { .* } as user
		`,
		map[string]interface{}{
			"passcode": leaguePasscode,
		})
		if err != nil {
			return []leagues.League{}, err
		}
		var us []users.User
		for result.Next() {
			record := result.Record()
			user, found := record.Get("user")
			if found {
				r := user.(map[string]interface{})
				u := users.User{
					Email: r["email"].(string),
				}
				us = append(us, u)
			}
		}
		return us, err
	})
	us, found := res.([]users.User)
	if !found {
		return []users.User{}, errors.New("Could not find the team in db results")
	}
	return us, err
}