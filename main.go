package main

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/controllers"
	"hilgardvr/ff1-go/service"
	"log"
	"net/http"
)


func main() {
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Init(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server on port", config.AppPort)
	http.HandleFunc(controllers.LoginCode, controllers.LoginCodeHandler)
	http.HandleFunc(controllers.Login, controllers.LoginHandler)
	http.HandleFunc(controllers.Logout, controllers.LoginHandler)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc(controllers.PickTeam, controllers.HomeContoller)
	http.HandleFunc(controllers.Home, controllers.HomeContoller)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
