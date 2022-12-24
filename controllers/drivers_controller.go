package controllers

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/helpers"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"log"
	"net/http"
)


func GetDrivers(w http.ResponseWriter, r *http.Request) {
	allDrivers, err := service.GetAllDrivers()
	if err != nil {
		log.Println("Unable to load drivers:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}

func GetBudget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(1000000)
}

func SaveController(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
	}
	var ds []drivers.Driver
	err = json.NewDecoder(r.Body).Decode(&ds)
	if err != nil {
		log.Println("Failed to parse drivers", err)
	}
	err = service.UpsertTeam(user, ds)
	if err != nil {
		log.Println("Failed to upsert team:", err)
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateLeagueController(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("Could not parse form")
		return
	}
	leagueName := r.Form.Get("league-name")
	generatedCode := helpers.GenerateLoginCode()
	err = service.SaveLeague(user, leagueName, generatedCode)
	if err != nil {
		log.Println("Could not save league:", err)
		return 
	}
	sub := fmt.Sprintf("League passcode for %s", leagueName)
	err = service.SendEmail(user.Email, sub, generatedCode)
	if err != nil {
		log.Println("Failed to send email: ", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}