package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
)

type Repo interface {
	Init(config *config.Config) error
	GetDrivers() []drivers.Driver
}