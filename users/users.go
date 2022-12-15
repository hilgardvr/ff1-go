package users

import "hilgardvr/ff1-go/drivers"

type User struct {
	Email string
	Team []drivers.Driver
}