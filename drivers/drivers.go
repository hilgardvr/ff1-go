package drivers

import (
	"hilgardvr/ff1-go/constructor"
	// "math"
)

type Driver struct {
	Id     	int64     	`json:"id"`
	Name   	string  	`json:"name"`
	Surname string  	`json:"surname"`
	Points 	int64     	`json:"points"`
	Price  	int64 		`json:"price"`
	Constructor constructor.Constructor `json:"constructor"`
}