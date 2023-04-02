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
	http.HandleFunc("/logincode", controllers.LoginCodeHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/league", controllers.LeagueController)
	http.HandleFunc("/pick-team", controllers.PickTeamController)
	http.HandleFunc("/repick-team", controllers.PickTeamController)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	http.HandleFunc("/api/save-team", controllers.SaveDriversController)
	http.HandleFunc("/api/create-league", controllers.CreateLeagueController)
	http.HandleFunc("/api/join-league", controllers.JoinLeagueController)
	http.HandleFunc("/api/team-details", controllers.SaveTeamDetails)
	http.HandleFunc("/display-league", controllers.DislayLeagueController)
	http.HandleFunc("/display-leagues", controllers.DisplayLeaguesController)
	http.HandleFunc("/display-member-team", controllers.DisplayTeamMemberController)
	http.HandleFunc("/display-points", controllers.DisplayRacePoints)
	http.HandleFunc("/admin/admin-page", controllers.CreateRacePoints)
	http.HandleFunc("/admin/update-data", controllers.UpdateRaceData)
	http.HandleFunc("/", controllers.HomeContoller)
	log.Println("Starting server on port", config.AppPort)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
