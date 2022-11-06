package repo

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/users"
	"net/http"

	// "io"
	// "io/ioutil"
	// "net/http"
	// "os/user"

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

func (n Neo4jRepo)AddUser(user users.User) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			merge (u:User {email: $email})
		`,
		map[string]interface{}{
			"email": user.Email,
		})
		return user.Email, err
	})
	return err
}

func (n Neo4jRepo)SaveSession(users.User, http.Cookie) error {
	return nil
}

func (n Neo4jRepo)GetUserFromSession(http.Cookie) (users.User, error) {
	return users.User{}, nil
}