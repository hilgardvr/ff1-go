package repo

import (
	"errors"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.Driver
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
	n.driver = driver
	return nil
}
	
func (n Neo4jRepo)GetDrivers() []drivers.Driver {
	return []drivers.Driver{}
}

func (n Neo4jRepo)AddUser(user users.User) (users.User, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			merge (u:User {Email: $email})
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		record, err := result.Single()
		if err != nil {
			return users.User{}, err
		}
		user, found := record.Get("user")
		if !found {
			return users.User{}, errors.New("Could not find user in result")
		}
		return user, nil
		// return user.Email, err
	})
	email, found := result.(map[string]interface{})["Email"] 
	if !found {
		return users.User{}, err
	}
	u := users.User{
		SessionId: "",
		Email: email.(string),
	}
	return u, err
}
