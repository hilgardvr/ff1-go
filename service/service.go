package service

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/repo"
	"hilgardvr/ff1-go/users"
	"log"
)

var svc ServiceIO

type ServiceIO struct {
	Db repo.Repo
	EmailService email.EmailService
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

func GetAllDrivers() ([]drivers.Driver, error) {
	allDrivers, err := svc.Db.GetDrivers()
	return allDrivers, err
}

func ValidateLoginCode(email string, code string) bool {
	return svc.Db.ValidateLoginCode(email, code)
}

func SendEmail(email string, subject string, body string) error {
	err := svc.EmailService.SendEmail(email, "Your F1-Go login code", body)
	return err
}

func SetLoginCode(email string, loginCode string) (string, error) {
	return svc.Db.SetLoginCode(email, loginCode)
}