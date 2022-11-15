package users

import "database/sql/driver"

type User struct {
	Email string
	Team []driver.Driver
}