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
 	// GetSession(r *http.Request) (users.User, error)
	SetLoginCode(email string) string
	DeleteLoginCode(email string)
	ValidateLoginCode(email string, codeToTest string) bool
}