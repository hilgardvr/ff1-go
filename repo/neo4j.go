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
			merge (u:User {email: $email})
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
	})
	email, found := result.(map[string]interface{})["email"] 
	if !found {
		return users.User{}, err
	}
	u := users.User{
		Email: email.(string),
	}
	return u, err
}

func (n Neo4jRepo) SetLoginCode(email string, generatedCode string) (string, error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			merge (u:User {email: $email})
			set u.loginCode = $generatedCode
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"generatedCode": generatedCode,
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
	})
	code, found := result.(map[string]interface{})["loginCode"] 
	if !found {
		return "", err
	}
	c := code.(string)
	return c, err
}

func (n Neo4jRepo)DeleteLoginCode(email string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			set u.loginCode = null
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
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
	})
	return err
}

func (n Neo4jRepo) ValidateLoginCode(email string, codeToTest string) bool {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email, loginCode: $loginCode})
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"loginCode": codeToTest,
		})
		record, err := result.Single()
		return record, err
	})
	return err == nil
}

func (n Neo4jRepo)SaveSession(email, uuid string) error {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {email: $email})
			set u.sessionId = $uuid
			return u { .* } as user
		`,
		map[string]interface{}{
			"email": email,
			"uuid": uuid,
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
	})
	return err
}

func (n Neo4jRepo)GetSession(uuid string) (string, bool) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		session.Close()
	}()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			match (u:User {sessionId: $uuid})
			return u { .* } as user
		`,
		map[string]interface{}{
			"uuid": uuid,
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
	})
	if err != nil {
		return "", false
	}
	email, found := result.(map[string]interface{})["email"] 
	return email.(string), found
}
