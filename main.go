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
	http.HandleFunc(controllers.LoginCode, controllers.LoginCodeHandler)
	http.HandleFunc(controllers.Login, controllers.LoginHandler)
	http.HandleFunc(controllers.Logout, controllers.LogoutHandler)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc("/api/budget", controllers.GetBudget)
	http.HandleFunc("/api/save-team", controllers.SaveController)
	http.HandleFunc(controllers.PickTeam, controllers.HomeContoller)
	http.HandleFunc(controllers.Home, controllers.HomeContoller)
	log.Println("Starting server on port", config.AppPort)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
