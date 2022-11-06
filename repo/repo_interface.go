package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"
	"net/http"
)

type Repo interface {
	Init(*config.Config) error
	GetDrivers() []drivers.Driver
	AddUser(users.User) error
	SaveSession(users.User, http.Cookie) error
	GetUserFromSession(http.Cookie) (users.User, error)
}