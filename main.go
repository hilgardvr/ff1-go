package main

import (
	"fmt"
	"hilgardvr/ff1-go/repo"
	"log"
)

func main() {
	err := repo.Init()
	if err != nil {
		log.Fatal(err)
	}
	drivers := repo.GetDrivers()
	fmt.Println(drivers)
}
