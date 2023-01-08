package service

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/races"
	"hilgardvr/ff1-go/repo"
	"hilgardvr/ff1-go/users"
	"log"
	"sort"
)

var svc ServiceIO

type ServiceIO struct {
	Db repo.Repo
	EmailService email.EmailService
	//todo make toggles
	SendEmail bool
}

func GetServiceIO() *ServiceIO {
	return &svc
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
		return err
	}
	svc = ServiceIO{
		Db: r,
		EmailService: e,
		SendEmail: config.SendEmails,
	}
	return nil
}


func UpsertTeam(user users.User, ds []drivers.Driver) error {
	valid := drivers.ValidateTeam(ds)
	if valid {
		err := svc.Db.DeleteTeam(user)
		if err != nil {
			log.Println("Failed to delete team for user: ", user, err)
			return err
		}
		err = svc.Db.SaveTeam(user, ds)
		if err != nil {
			log.Println("Failed to save team for user: ", user, err)
			return err
		}
	}
	return nil
}

func GetUserTeam(user users.User) (users.User, error) {
	team, err := svc.Db.GetTeam(user)
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

func GetAllRaces() ([]races.Race, error) {
	allRaces, err := svc.Db.GetAllRaces()
	return allRaces, err
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

func GetLeagueUsers(passcode string) ([]users.User, error) {
	return svc.Db.GetLeagueMembers(passcode)
}

func CreateRacePoints(racePoints []drivers.Driver) error {
	race, err := GetLatestRace()
	if err != nil {
		log.Println("Error getting latest race:", err)
		return err
	}
	race.Race += 1
	return svc.Db.CreateNewRace(racePoints, race)
}