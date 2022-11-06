package service

import (
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/repo"
	"hilgardvr/ff1-go/config"
)

var svc ServiceIO

type ServiceIO struct {
	Db repo.Repo
	EmailService email.EmailService
}

func GetService() *ServiceIO {
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