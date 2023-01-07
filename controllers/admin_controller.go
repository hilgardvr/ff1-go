package controllers

import (
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)

var season = 2022

func CreateRacePoints(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
		return
	}
	if user.IsAdmin {
		drivers, err := service.GetAllDriversForSeason(season)
		if err != nil {
			log.Println("error fetching drivers:", err)
			return
		}
		view.AdminTemplate(w, user, season, drivers)
	} else {
		http.Redirect(w, r, "/", 300)
	}
}