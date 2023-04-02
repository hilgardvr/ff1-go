package repo

import (
	"errors"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/constructor"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/races"
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

func (n Neo4jRepo)GetDriversBySeason(season int) ([]drivers.Driver,  error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (d:Driver)-[hr:HAS_RACE]->(r:Race {season: $season}) 
			match (d)-[rf:RACES_FOR]->(c:Constructor)
			return ID(d) as id, d.name as name, d.surname as surname, sum(hr.points) as points, ID(c) as constructorId, c.name as constructorName, d.price as price
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
			constructorId, found := record.Get("constructorId")
			if !found {
				log.Println("Could not find constructor id")
				continue
			}
			constructorName, found := record.Get("constructorName")
			if !found {
				log.Println("Could not find constructor name")
				continue
			}
			price, found := record.Get("price")
			if !found {
				log.Println("Could not find driver price")
				continue
			}
			if price == nil {
				price = int64(1000000)
			}
			driver := drivers.Driver{
				Id: id.(int64),
				Name: name.(string),
				Surname: surname.(string),
				Points: points.(int64),
				Constructor: constructor.Constructor{
					Id: constructorId.(int64),
					ConstructorName: constructorName.(string),
				},
				Price: price.(int64),
			}
			ls = append(ls, driver)
		}

		return ls, err
	})
	if err != nil {
		return []drivers.Driver{}, err
	}
	ds := result.([]drivers.Driver)
	return ds, err
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
			with u,
			case
				when u.budget is null then 1000000
				else u.budget
			end as b
			set u.budget = b
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
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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
	res := result.(map[string]interface{})
	email, found := res["email"] 
	if !found {
		return users.User{}, false
	}
	isAdmin, found := res["isadmin"] 
	if !found {
		isAdmin = false
	}
	teamName, found := res["teamName"] 
	if !found {
		teamName = ""
	}
	teamPriciple, found := res["teamPrinciple"] 
	if !found {
		teamPriciple = ""
	}
	budget, found := res["budget"] 
	if !found {
		log.Println("No budget found for user: ", email)
		budget = 0
	}
	n.DeleteLoginCode(email.(string))
	return users.User{
		Email: email.(string), 
		IsAdmin: isAdmin.(bool),
		TeamName: teamName.(string),
		TeamPriciple: teamPriciple.(string),
		Budget: budget.(int64),
	}, true
}

func (n Neo4jRepo) UpdateTeam(user users.User, selectedDrivers []drivers.Driver, race races.Race) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	var driverIds []int64
	for _, d := range selectedDrivers {
		driverIds = append(driverIds, d.Id)
	}
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})-[:HAS_TEAM]->(t:Team)-[hd:HAS_DRIVER]->(:Driver)
			where (t)-[:FOR_RACE]->(:Race {season: $season, race: $race})
			detach delete t
		`,
		map[string]interface{}{
			"email": user.Email,
			"season": race.Season,
			"race": race.Race,
		})
		if err != nil {
			log.Println("Error deleting updated team team:", err)
			return nil, err 
		}
		_, err = result.Consume()
		if err != nil {
			log.Println("Error consuming update team:", err)
			return nil, err 
		}
		result, err = tx.Run(`
			match (u:User {email: $email})
			match (r:Race {season: $season, race: $race})
			match (d:Driver) where ID(d) in $driverIds
			merge (t:Team {email: $email})-[:FOR_RACE]->(r)
			merge (u)-[h:HAS_TEAM]->(t)
			merge (t)-[hd:HAS_DRIVER]->(d)
			set u.budget = $newBudget
			return u { .* } as user
		`, 
		map[string]interface{}{
			"email": user.Email,
			"season": race.Season,
			"race": race.Race,
			"driverIds": driverIds,
			"newBudget": user.Budget,
		})
		if err != nil {
			log.Println("Error deleting updated team team:", err)
			return nil, err 
		}
		_, err = result.Consume()
		if err != nil {
			log.Println("Error consuming update team:", err)
			return nil, err 
		}
		return  nil, err
	})
	return err
}

func (n Neo4jRepo) AddUsersRaceBudget(additionalAmount int) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			match (u:User)
			set u.budget = u.budget + $addAmt
			return u
		`,
		map[string]interface{}{
			"addAmt": additionalAmount,
		})
		return nil, err
	})
	if err != nil {
		log.Println("Error updating user budgets: ", err)
	}
	return err
}

// func (n Neo4jRepo) DeleteTeam(user users.User, race races.Race) error {
// 	session := n.driver.NewSession(neo4j.SessionConfig{})
// 	defer func() {
// 		session.Close()
// 	}()
// 	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
// 		result, err := tx.Run(`
// 			match (u:User {email: $email})-[:HAS_TEAM]->(t:Team)
// 			where (t)-[:FOR_RACE]->(:Race {season: $season, race: $race})
// 			detach delete t
// 		`,
// 		map[string]interface{}{
// 			"email": user.Email,
// 			"season": race.Season,
// 			"race": race.Race,
// 		})
// 		resultSummary, err := result.Consume()
// 		if err != nil {
// 			return []drivers.Driver{}, err
// 		}
// 		return resultSummary, err
// 	})
// 	return err
// }

// func (n Neo4jRepo) SaveTeam(user users.User, selectedDrivers []drivers.Driver, race races.Race) error {
// 	session := n.driver.NewSession(neo4j.SessionConfig{})
// 	defer func() {
// 		session.Close()
// 	}()
// 	var driverIds []int64
// 	for _, d := range selectedDrivers {
// 		driverIds = append(driverIds, d.Id)
// 	}
// 	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
// 		result, err := tx.Run(`
// 			match (u:User {email: $email})
// 			match (r:Race {season: $season, race: $race})
// 			match (d:Driver) where ID(d) in $driverIds
// 			merge (t:Team {email: $email})-[:FOR_RACE]->(r)
// 			merge (u)-[h:HAS_TEAM]->(t)
// 			merge (t)-[hd:HAS_DRIVER]->(d)
// 			return u { .* } as user
// 		`,
// 		map[string]interface{}{
// 			"email": user.Email,
// 			"season": race.Season,
// 			"race": race.Race,
// 			"driverIds": driverIds,
// 		})
// 		if err != nil {
// 			log.Println("Error save team:", err)
// 			return users.User{}, errors.New("Error running save team query")
// 		}
// 		_, err = result.Consume()
// 		if err != nil {
// 			log.Println("Error creating team", err)
// 			return users.User{}, err
// 		}
// 		return users.User{}, nil
// 	})
// 	return err
// }

func (n Neo4jRepo) GetUserDetails(email string) (users.User, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
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
	if err != nil {
		return users.User{}, err
	}
	res := result.(map[string]interface{})
	isAdmin, found := res["isadmin"] 
	if !found {
		isAdmin = false
	}
	teamName, found := res["teamName"] 
	if !found {
		teamName = ""
	}
	teamPriciple, found := res["teamPrinciple"] 
	if !found {
		teamPriciple = ""
	}
	budget, found := res["budget"] 
	if !found {
		log.Println("No budget found for user: ", email)
		budget = 0
	}
	return users.User{
		Email: email, 
		IsAdmin: isAdmin.(bool),
		TeamName: teamName.(string),
		TeamPriciple: teamPriciple.(string),
		Budget: budget.(int64),
	}, nil
}

func (n Neo4jRepo) GetUserTeamForRace(user users.User, race races.Race) ([]drivers.Driver, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	allDrivers, err := n.GetDriversBySeason(int(race.Season))
	if err != nil {
		log.Println("Could not get all drivers by season:", err)
		return []drivers.Driver{}, err
	}
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})-[:HAS_TEAM]->(t:Team)-[:HAS_DRIVER]-(d:Driver)
			where (t)-[:FOR_RACE]->(:Race {season: $season, race: $race})
			optional match (d)-[hr:HAS_RACE]-(:Race)
			return ID(d) as id, sum(hr.points) as points
		`,
		map[string]interface{}{
			"email": user.Email,
			"season": race.Season,
			"race": race.Race,
		})
		var ls []drivers.Driver
		for result.Next() {
			record := result.Record()
			id, found := record.Get("id")
			if !found {
				log.Println("Could not find driver id")
				continue
			}
			for _, v := range allDrivers {
				if v.Id == id.(int64) {
					ls = append(ls, v)
				}
			}
		}

		return ls, err
	})
	if err != nil {
		return []drivers.Driver{}, err
	}
	ds := result.([]drivers.Driver)
	return ds, err
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

func (n Neo4jRepo) GetRacePoints(race races.Race) (races.RacePoints, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (d:Driver)
			optional match (d)-[hr:HAS_RACE]-(r:Race {season: $season, race: $race})
			return d {.*, driverId: ID(d)} as driver, hr.points as driverPoints
		`,
		map[string]interface{}{
			"season": race.Season,
			"race": race.Race,
		})
		if err != nil {
			return []leagues.League{}, err
		}
		var dr []drivers.Driver
		for result.Next() {
			record := result.Record()
			var points int64
			pointsRec, found := record.Get("driverPoints")
			if found && pointsRec != nil {
				points = pointsRec.(int64)
			} else {
				points = 0
			}
			driver, found := record.Get("driver")
			if found {
				rec := driver.(map[string]interface{})
				d := drivers.Driver{
					Id: rec["driverId"].(int64),
					Name: rec["name"].(string),
					Surname: rec["surname"].(string),
					Points: points,
				}
				dr = append(dr, d)
			}
		}
		return dr, err
	})
	ds, found := res.([]drivers.Driver)
	if !found {
		return races.RacePoints{}, errors.New("Could not find the drivers in db results")
	}
	return races.RacePoints{Race: race, Drivers: ds}, err
}

func (n Neo4jRepo) GetLeagueMembers(leaguePasscode string, season int) ([]users.User, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (l:League {passcode: $passcode})-[:LEAGUE]-(u:User)
			optional match (u)-[ht:HAS_TEAM]-(t:Team)-[fr:FOR_RACE]-(r:Race {season: $season})
			optional match (d:Driver)-[hr:HAS_RACE]-(r)
			where (t)-[:HAS_DRIVER]-(d)
			with sum(hr.points) as p, u as u
			return u {. *, points: p} as user
		`,
		map[string]interface{}{
			"passcode": leaguePasscode,
			"season": season,
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
				email := r["email"].(string)
				teamName := ""
				if r["teamName"] != nil {
					teamName = r["teamName"].(string)
				}
				teamPrinciple := ""
				if r["teamPrinciple"] != nil {
					teamPrinciple = r["teamPrinciple"].(string)
				}
				u := users.User{
					Email: email,
					SeasonPoints: int(r["points"].(int64)),
					TeamName: teamName,
					TeamPriciple: teamPrinciple,
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

func (n Neo4jRepo) GetAllCompletedRaces() ([]races.Race, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (:Driver)-[:HAS_RACE]-(r:Race) return distinct r { .* } as race
		`, map[string]interface{}{})
		if err != nil {
			return []races.Race{}, err
		}
		var rs []races.Race
		for result.Next() {
			record := result.Record()
			race, found := record.Get("race")
			if found {
				r := race.(map[string]interface{})
				t, ok := r["track"]
				track := ""
				if ok {
					track = t.(string)
				}
				u := races.Race{
					Race: r["race"].(int64),
					Season: r["season"].(int64),
					Track: track,
				}
				rs = append(rs, u)
			}
		}
		return rs, err
	})
	rs, found := res.([]races.Race)
	if !found {
		return []races.Race{}, errors.New("Could not find the team in db results")
	}
	return rs, err
}

func (n Neo4jRepo) GetAllRaces() ([]races.Race, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (r:Race) return r { .* } as race
		`, map[string]interface{}{})
		if err != nil {
			return []races.Race{}, err
		}
		var rs []races.Race
		for result.Next() {
			record := result.Record()
			race, found := record.Get("race")
			if found {
				r := race.(map[string]interface{})
				t, ok := r["track"]
				track := ""
				if ok {
					track = t.(string)
				}
				u := races.Race{
					Race: r["race"].(int64),
					Season: r["season"].(int64),
					Track: track,
				}
				rs = append(rs, u)
			}
		}
		return rs, err
	})
	rs, found := res.([]races.Race)
	if !found {
		return []races.Race{}, errors.New("Could not find the races")
	}
	return rs, err
}

func (n Neo4jRepo) CreateNewRace(driverWithPoints []drivers.Driver, race races.Race, track string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			merge (r:Race {season: $season, race: $race})
			return r
		`,
		map[string]interface{}{
			"season": race.Season,
			"race": race.Race + 1,
		})
		if err != nil {
			return "", err
		}
		for _, v := range driverWithPoints {
			_, err := tx.Run(`
				match (r:Race {season: $season, race: $race})
				set r.track = $track
				with r as r
				match (d:Driver)
				where ID(d) = $id
				merge (d)-[:HAS_RACE {points: $points}]-(r)
			`,
			map[string]interface{}{
				"season": race.Season,
				"race": race.Race,
				"id": v.Id,
				"points": v.Points,
				"track": track,
			})
			if err != nil {
				log.Println("Error adding drivers to race:", err)
				return "", err
			}
		}
		return "", nil
	})
	return err
}

func (n Neo4jRepo) SaveUserTeamDetails(user users.User) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			set u.teamName = $teamName
			set u.teamPrinciple = $teamPrinciple
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": user.Email,
			"teamName": user.TeamName,
			"teamPrinciple": user.TeamPriciple,
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


func (n Neo4jRepo)SetDriverPrice(pricedDriver drivers.Driver) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			match (d:Driver)
			where ID(d) = $did
			set d.price = $price
		`,
		map[string]interface{}{
			"did": pricedDriver.Id,
			"price": pricedDriver.Price,
		})
		if err != nil {
			return drivers.Driver{}, err
		}
		return drivers.Driver{}, nil
	})
	return err
}