package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/users"
	"time"
)

type Repo interface {
	Init(*config.Config) error
	GetDrivers() ([]drivers.Driver, error)
	AddUser(users.User) (users.User, error)

	SetLoginCode(email string, generatedCode string) (string, error)
	DeleteLoginCode(email string) error
	ValidateLoginCode(email string, codeToTest string) bool

	SaveSession(email, uuid string, duration time.Duration) error
	GetSession(uuid string) (string, bool)

	SaveTeam(users.User, []drivers.Driver) error
	GetTeam(users.User) ([]drivers.Driver, error)
	DeleteTeam(users.User) error

	SaveLeague(user users.User, leagueName string, passcode string) error
	GetLeagueForUser(user users.User) (leagues []leagues.League, err error)
	JoinLeague(user users.User, leaguePasscode string) error
}