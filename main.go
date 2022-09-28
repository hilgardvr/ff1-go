package main

import (
	"fmt"
	"hilgardvr/ff1-go/controllers"
	"hilgardvr/ff1-go/repo"
	"log"
	"net/http"
)

const port = ":9000"

func main() {
	err := repo.Init()
	if err != nil {
		log.Fatal(err)
	}
	drivers := repo.GetDrivers()
	fmt.Println(drivers)
	http.HandleFunc("/", controllers.HomeController)
	log.Fatal(http.ListenAndServe(port, nil))
}
