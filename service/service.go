package service

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/pricing"
	"hilgardvr/ff1-go/races"
	"hilgardvr/ff1-go/repo"
	"hilgardvr/ff1-go/users"
	"log"
	"sort"
)

var svc ServiceIO
const budget float64 = 1000000.0
const driversInTeam int = 4
const basePrice float64 = budget * 0.1
const adjustmentFactor float64 = 1.5


type ServiceIO struct {
	Db repo.Repo
	EmailService email.EmailService
	SendEmail bool
}

func GetServiceIO() *ServiceIO {
	return &svc
}

func AssingDriverPrices() error {
	d, err := GetAllDriverForCurrentSeason()
	if err != nil {
		log.Println("could not get drivers for current season", err)
		return err
	}
	pricedDrivers := pricing.AssignPrices(d)
	for _, v := range pricedDrivers {
		err = svc.Db.SetDriverPrice(v)
		if err != nil {
			log.Println("could not set driver price", err)
			return err
		}
	}
	return  nil
}

func Init(config *config.Config) error {
	r := &repo.Neo4jRepo{}
	err := r.Init(config)
	if err != nil {
		return err
	}
	e := &email.GmailEmailService{}
	err = e.Init(config)
	if err != nil {
		log.Println("Error Init:", err)
		return err
	}
	svc = ServiceIO{
		Db: r,
		EmailService: e,
		SendEmail: config.SendEmails,
	}
	err = AssingDriverPrices()
	if err != nil {
		log.Println("Error assigning prices:", err)
		return err
	}
	return nil
}


func UpsertTeam(user users.User, ds []drivers.Driver) error {
	valid := drivers.ValidateTeam(ds, user.Budget)
	if valid {
		latesstRace, err := GetLatestRace()
		if err != nil {
			return nil
		}
		err = svc.Db.DeleteTeam(user, latesstRace)
		if err != nil {
			log.Println("Failed to delete team for user: ", user, err)
			return err
		}
		err = svc.Db.SaveTeam(user, ds, latesstRace)
		if err != nil {
			log.Println("Failed to save team for user: ", user, err)
			return err
		}
	}
	return nil
}

func GetUserDetails(email string) (users.User, error) {
	user, err := svc.Db.GetUserDetails(email)
	if err != nil {
		log.Println("Error looking up user details", err)
	}
	return user, err
}

func GetUserTeam(user users.User) (users.User, error) {
	latestRace, err := GetLatestRace()
	if err != nil {
		log.Println("Failed to fetch latest race: ", user, err)
		return users.User{}, err
	}
	team, err := svc.Db.GetUserTeamForRace(user, latestRace)
	if err != nil {
		log.Println("Failed to fetch user team for user: ", user, err)
		return users.User{}, err
	}
	user = users.User{
			Email: user.Email,
			Team: team,
		}
	return user, nil
}

func GetAllDriverForCurrentSeason() ([]drivers.Driver, error) {
	latestRace, err := GetLatestRace()
	if err != nil {
		log.Println("Could not get latest race", err)
		return []drivers.Driver{}, err
	}
	return GetAllDriversForSeason(int(latestRace.Season))
}

func GetAllDriversForSeason(season int) ([]drivers.Driver, error) {
	allDrivers, err := svc.Db.GetDriversBySeason(season)
	if err != nil {
		return []drivers.Driver{}, err
	}
	sort.Slice(allDrivers, func(i, j int) bool {
		return allDrivers[i].Points > allDrivers[j].Points
	})
	return allDrivers, err
}

func GetLatestRace() (races.Race, error) {
	allRaces, err := GetAllRaces()
	if err != nil {
		return races.Race{}, err
	}
	sort.Slice(allRaces, func(i, j int) bool {
		if allRaces[i].Season > allRaces[j].Season {
			return true
		} else {
			if allRaces[i].Season == allRaces[j].Season {
				return allRaces[i].Race > allRaces[j].Race
			} else {
				return false
			}
		}
	})
	return allRaces[0], nil
}

func GetAllCompletedRaces() ([]races.Race, error) {
	allCompletedRaces, err := svc.Db.GetAllCompletedRaces()
	return allCompletedRaces, err
}

func GetLatestCompletedRace() (races.Race, error) {
	allCompleted, err := GetAllCompletedRaces()
	if err != nil {
		log.Println("could not get all completed races")
		return races.Race{}, err
	}
	sort.Slice(allCompleted, func(i, j int) bool {
		if allCompleted[i].Season > allCompleted[j].Season {
			return true
		} else {
			if allCompleted[i].Season == allCompleted[j].Season {
				return allCompleted[i].Race > allCompleted[j].Race
			} else {
				return false
			}
		}
	})
	return allCompleted[0], nil
}

func SortRacesDesc(races []races.Race) (race []races.Race) {
	sort.Slice(races, func(i, j int) bool {
		if races[i].Season > races[j].Season {
			return true
		} else {
			if races[i].Season == races[j].Season {
				return races[i].Race > races[j].Race
			} else {
				return false
			}
		}
	})
	return races
}

func GetAllRaces() ([]races.Race, error) {
	allRaces, err := svc.Db.GetAllRaces()
	return allRaces, err
}

func GetAllRacesForCurrentSeason() ([]races.Race, error) {
	l, err := GetLatestRace()
	if err != nil {
		return []races.Race{}, err
	}
	return GetAllRacesForSeason(l.Season)
}

func GetAllCompletedRacesForCurrentSeason() ([]races.Race, error) {
	latestRace, err := GetLatestRace()
	if err != nil {
		return []races.Race{}, err
	}
	allCompleted, err := GetAllCompletedRaces()
	if err != nil {
		return []races.Race{}, err
	}
	completedCurrentSeason := []races.Race{}
	for _, v := range allCompleted {
		if v.Season == latestRace.Season {
			completedCurrentSeason = append(completedCurrentSeason, v)
		}
	}
	return completedCurrentSeason, nil
	
}

func GetAllRacesForSeason(season int64) ([]races.Race, error) {
	allRaces, err := GetAllRaces()
	if err != nil {
		return allRaces, err
	}
	var seasonRaces []races.Race
	for _, v := range allRaces {
		if v.Season == season {
			seasonRaces = append(seasonRaces, v)
		}
	}
	return SortRacesDesc(seasonRaces), err
}

func ValidateLoginCode(email string, code string) bool {
	return svc.Db.ValidateLoginCode(email, code)
}

func SendEmail(email string, subject string, body string) error {
	var err error
	if svc.SendEmail {
		err = svc.EmailService.SendEmail(email, subject, body)
	} else {
		log.Println("Sending of emails toggled off")
	}
	return err
}

func SetLoginCode(email string, loginCode string) (string, error) {
	return svc.Db.SetLoginCode(email, loginCode)
}

func SaveLeague(user users.User, leagueName string, passcode string) error {
	err := svc.Db.SaveLeague(user, leagueName, passcode)
	return err
}

func JoinLeague(user users.User, passcode string) error {
	err := svc.Db.JoinLeague(user, passcode)
	return err
}

func SaveUserTeamDetails(user users.User) error {
	return svc.Db.SaveUserTeamDetails(user)
}

func GetLeagueUsers(passcode string) ([]users.User, error) {
	latestRace, err := GetLatestRace()
	if err != nil {
		log.Println("Could not get latest race:", err)
		return []users.User{}, err
	}
	return svc.Db.GetLeagueMembers(passcode, int(latestRace.Season))
}

func CreateRacePoints(racePoints []drivers.Driver, track string) error {
	race, err := GetLatestRace()
	if err != nil {
		log.Println("Error getting latest race:", err)
		return err
	}
	return svc.Db.CreateNewRace(racePoints, race, track)
}

func GetUserRacePoints(user users.User, race races.Race) (races.RacePoints, error) {
	ut, err := svc.Db.GetUserTeamForRace(user, race)
	if err != nil {
		log.Println("Error getting user team for race:", err)
		return races.RacePoints{}, err
	}
	rp, err := svc.Db.GetRacePoints(race)
	if err != nil {
		log.Println("Error getting latest race:", err)
		return races.RacePoints{}, err
	}
	var teamWithPoints []drivers.Driver
	var total int64
	for _, v := range ut {
		for _, r := range rp.Drivers {
			if v.Id == r.Id {
				teamWithPoints = append(teamWithPoints, r)
				total += r.Points
			}
		}
	}
	return races.RacePoints{Race: race, Drivers: teamWithPoints, Total: total}, err
}

// func AssignPrices(drivers []drivers.Driver, currentRaces []races.Race) []drivers.Driver {
// 	var createdDrivers []drivers.Driver
// 	totalPoints := sumAllDriverPoints(drivers)
// 	for _, driver := range drivers {
// 		price := calcPrice(driver, totalPoints, len(currentRaces))
// 		createdDrivers = append(createdDrivers, drivers.Driver{
// 			Id:     driver.Id,
// 			Name:   driver.Name,
// 			Surname: driver.Surname,
// 			Points: driver.Points,
// 			Price:  price,
// 			Constructor: driver.Constructor,
// 		})
// 	}
// 	return createdDrivers
// }

// func sumAllDriverPoints(drivers []drivers.Driver) int64 {
// 	var totalPoints int64
// 	for _, driver := range drivers {
// 		totalPoints += driver.Points
// 	}
// 	return totalPoints
// }

// func calcPrice(driver drivers.Driver, totalPoints int64, numberOfRaces int) int64 {
// 	if totalPoints == 0 {
// 		totalPoints++
// 	}
// 	// driverPointsShare := float64(driver.Points) / float64(totalPoints)
// 	// price := (driverPointsShare*budget + basePrice) * adjustmentFactor
// 	driverPointsShare := float64(driver.Points) / float64(totalPoints)
// 	price := budget * driverPointsShare
// 	price = math.Round(price)
// 	return int64(price)
// }