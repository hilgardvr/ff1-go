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
const League = "/league"
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
	http.HandleFunc(League, controllers.LeagueController)
	http.HandleFunc(PickTeam, controllers.PickTeamController)
	http.HandleFunc(RepickTeam, controllers.PickTeamController)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc("/api/save-team", controllers.SaveDriversController)
	http.HandleFunc("/api/create-league", controllers.CreateLeagueController)
	http.HandleFunc("/api/join-league", controllers.JoinLeagueController)
	http.HandleFunc("/display-league", controllers.DislayLeagueController)
	http.HandleFunc("/admin/admin-page", controllers.CreateRacePoints)
	http.HandleFunc("/admin/update-data", controllers.UpdateRaceData)
	http.HandleFunc(Home, controllers.HomeContoller)
	log.Println("Starting server on port", config.AppPort)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
