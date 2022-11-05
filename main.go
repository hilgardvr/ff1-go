package main

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/controllers"
	"hilgardvr/ff1-go/service"
	"log"
	"net/http"
)

const port = ":3000"

func main() {
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Init(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server on port", port)
	http.HandleFunc("/logincode", controllers.LoginCodeHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc("/", controllers.HomeContoller)
	log.Fatal(http.ListenAndServe(port, nil))
}
