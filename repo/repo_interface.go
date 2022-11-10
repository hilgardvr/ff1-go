package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"
)

type Repo interface {
	Init(*config.Config) error
	GetDrivers() []drivers.Driver
	AddUser(users.User) (users.User, error)

	SetLoginCode(email string, generatedCode string) (string, error)
	DeleteLoginCode(email string) error
	ValidateLoginCode(email string, codeToTest string) bool

	SaveSession(email, uuid string) error
	GetSession(uuid string) (string, bool)
}