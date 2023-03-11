package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/leagues"
	"hilgardvr/ff1-go/races"
	"hilgardvr/ff1-go/users"
	"time"
)

type Repo interface {
	Init(*config.Config) error

	GetDriversBySeason(int) ([]drivers.Driver, error)

	SetLoginCode(email string, generatedCode string) (string, error)
	DeleteLoginCode(email string) error
	ValidateLoginCode(email string, codeToTest string) bool

	SaveSession(email, uuid string, duration time.Duration) error
	GetSession(uuid string) (users.User, bool)

	SaveTeam(users.User, []drivers.Driver, races.Race) error
	GetUserTeamForRace(user users.User, race races.Race) ([]drivers.Driver, error)
	DeleteTeam(users.User, races.Race) error
	SaveUserTeamDetails(users.User) error
	GetUserDetails(email string) (users.User, error)

	SaveLeague(user users.User, leagueName string, passcode string) error
	GetLeagueForUser(user users.User) (leagues []leagues.League, err error)
	JoinLeague(user users.User, leaguePasscode string) error
	GetLeagueMembers(leaguePasscode string, season int) ([]users.User, error)

	GetAllRaces() ([]races.Race, error)
	GetAllCompletedRaces() ([]races.Race, error)
	CreateNewRace(drivers []drivers.Driver, race races.Race, track string) error

	GetRacePoints(races.Race) (races.RacePoints, error)

	SetDriverPrice(drivers drivers.Driver) error
}