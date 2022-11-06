package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepo struct {
	driver *neo4j.Driver
}

func (n *Neo4jRepo)Init(config *config.Config) error {
	driver, err := neo4j.NewDriver(
		config.Neo4jUri,
		neo4j.BasicAuth(
			config.Neo4jUsername,
			config.Neo4jPassword,
			"",
		),
	)
	if err != nil {
		return err
	}
	err = driver.VerifyConnectivity()

	if err != nil {
		return err
	}
	n.driver = &driver
	return nil
}
	
func (n Neo4jRepo)GetDrivers() []drivers.Driver {
	return []drivers.Driver{}
}