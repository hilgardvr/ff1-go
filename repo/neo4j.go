package repo

import (
	"errors"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/users"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.Driver
}

const migrationFilePath = "./repo/data"

func migrate(d neo4j.Driver) error {
	fs, err := ioutil.ReadDir(migrationFilePath)
	if err != nil {
		log.Println("Error opening directory:", err)
		return err
	}
	var cypherMigrtions []string
	for _, n := range fs {
		if strings.Contains(n.Name(), ".cypher") {
			cypherMigrtions = append(cypherMigrtions, migrationFilePath + "/" + n.Name())
		}
	}
	sort.Slice(cypherMigrtions, func(i, j int) bool {
		is := filepath.Base(strings.Split(cypherMigrtions[i], "_")[0])
		js := filepath.Base(strings.Split(cypherMigrtions[j], "_")[0])
		ii, err := strconv.ParseInt(is, 10, 64)
		if err != nil {
			log.Println("Could not parse migration timestamp:", err)
		}
		ij, err := strconv.ParseInt(js, 10, 64)
		if err != nil {
			log.Println("Could not parse migration timestamp:", err)
		}
		return  ii < ij
	})
	executeMigrations(d, cypherMigrtions)
	return nil
}

func executeMigrations(d neo4j.Driver, paths []string) error {
	session := d.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		for _, p := range paths {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				log.Println("Unable to read migration file:", err)
				return nil, err
			}
			_, err = tx.Run(string(b), map[string]interface{}{})
			if err != nil {
				log.Println("Error migrating:", p, err)
				return nil, err
			}
			log.Println("Migration successful for: ", p)
		}
		return nil, nil
	})
	return err
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
	err = migrate(driver)
	if err != nil {
		return err
	}
	n.driver = driver
	return err
}

	
func (n Neo4jRepo)GetDriversBySeason(season int) (neo4jDriver []drivers.Driver, err error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (d:Driver)-[hr:HAS_RACE]->(r:Race {season: $season}) 
			return ID(d) as id, d.name as name, d.surname as surname, sum(hr.points) as points
		`,
		map[string]interface{}{
			"season": season,
		})
		if err != nil {
			return []drivers.Driver{}, err
		}
		var ls []drivers.Driver
		for result.Next() {
			record := result.Record()
			id, found := record.Get("id")
			if !found {
				log.Println("Could not find driver id")
				continue
			}
			name, found := record.Get("name")
			if !found {
				log.Println("Could not find driver name")
				continue
			}
			surname, found := record.Get("surname")
			if !found {
				log.Println("Could not find driver surname")
				continue
			}
			points, found := record.Get("points")
			if !found {
				log.Println("Could not find driver points")
				continue
			}
			driver := drivers.Driver{
				Id: id.(int64),
				Name: name.(string),
				Surname: surname.(string),
				Points: points.(int64),
			}
			ls = append(ls, driver)
		}

		return ls, err
	})
	if err != nil {
		return []drivers.Driver{}, err
	}
	ds := result.([]drivers.Driver)
	ds = drivers.AssignPrices(ds)
	return ds, err
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

func (n Neo4jRepo) GetSession(uuid string) (users.User, bool) {
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
		return users.User{}, false
	}
	email, found := result.(map[string]interface{})["email"] 
	if !found {
		return users.User{}, false
	}
	isAdmin, found := result.(map[string]interface{})["isadmin"] 
	if !found {
		isAdmin = false
		found = true
	}
	n.DeleteLoginCode(email.(string))
	return users.User{Email: email.(string), IsAdmin: isAdmin.(bool)}, found
}

func (n Neo4jRepo) SaveTeam(user users.User, selectedDrivers []drivers.Driver) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	var drivers []int64
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
	var ids []int64
	for _, d := range res {
		ids = append(ids, d.(int64))
	}
	allDrivers, err := n.GetDriversBySeason(2022)
	if err != nil {
		log.Println("Could not get drivers by season:", err)
		return []drivers.Driver{}, err
	}
	var ds []drivers.Driver
	for _, ad := range allDrivers {
		for _, ud := range ids {
			if ud == ad.Id {
				ds = append(ds, ad)
			}
		}
	}
	return ds, err
}

// func (n Neo4jRepo) GetDriversByIdForSeason(ids []int64, season int) ([]drivers.Driver, error) {
// 	session := n.driver.NewSession(neo4j.SessionConfig{})
// 	defer func() {
// 		session.Close()
// 	}()
// 	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
// 		result, err := tx.Run(`
// 			match (d:Driver)-[hr:HAS_RACE]-(r:Race {season: $season})
// 			where ID(d) in $ids
// 			return ID(d) as id, d.name as name, d.surname as surname, sum(hr.points) as points
// 		`,
// 		map[string]interface{}{
// 			"ids": ids,
// 			"season": season,
// 		})
// 		if err != nil {
// 			return []drivers.Driver{}, err
// 		}
// 		var ds []drivers.Driver
// 		for result.Next() {
// 			record := result.Record()
// 			id, found := record.Get("id")
// 			if !found {
// 				log.Println("Could not find name")
// 				return []drivers.Driver{}, errors.New("Could not extract name")
// 			}
// 			name, found := record.Get("name")
// 			if !found {
// 				log.Println("Could not find name")
// 				return []drivers.Driver{}, errors.New("Could not extract name")
// 			}
// 			surname, found := record.Get("surname")
// 			if !found {
// 				log.Println("Could not find name")
// 				return []drivers.Driver{}, errors.New("Could not extract name")
// 			}
// 			points, found := record.Get("points")
// 			if !found {
// 				log.Println("Could not find name")
// 				return []drivers.Driver{}, errors.New("Could not extract name")
// 			}
// 			l := drivers.Driver{
// 				Id: id.(int64),
// 				Name: name.(string),
// 				Surname: surname.(string),
// 				Points: points.(int64),
// 			}
// 			ds = append(ds, l)
// 		}
// 		return ds, err
// 	})
// 	if err != nil {
// 		return []drivers.Driver{}, err
// 	}
// 	parsedDrivers, found := result.([]drivers.Driver)
// 	if !found {
// 		return []drivers.Driver{}, errors.New("Could not find the drivers in db results")
// 	}
// 	return parsedDrivers, err
// }

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
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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
	if err != nil {
		return []leagues.League{}, err
	}
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