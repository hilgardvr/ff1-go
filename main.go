package main

import (
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/controllers"
	"hilgardvr/ff1-go/service"
	"log"
	"net/http"
)

const PickTeam = "/pick-team"
const RepickTeam = "/repick-team"
const LoginCode = "/logincode"
const Login = "/login"
const Logout = "/logout"
const Home = "/"


func main() {
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Init(config)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc(LoginCode, controllers.LoginCodeHandler)
	http.HandleFunc(Login, controllers.LoginHandler)
	http.HandleFunc(Logout, controllers.LogoutHandler)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc("/api/budget", controllers.GetBudget)
	http.HandleFunc("/api/save-team", controllers.SaveController)
	http.HandleFunc(PickTeam, controllers.PickTeamController)
	http.HandleFunc(RepickTeam, controllers.PickTeamController)
	http.HandleFunc(Home, controllers.HomeContoller)
	log.Println("Starting server on port", config.AppPort)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
