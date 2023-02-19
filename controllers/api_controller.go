package controllers

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/helpers"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/users"
	"log"
	"net/http"
	"strings"
)

func GetDrivers(w http.ResponseWriter, r *http.Request) {
	latestRace, err := service.GetLatestRace()
	if err != nil {
		log.Println("Unable to get latest race:", err)
		return
	}
	allDrivers, err := service.GetAllDriversForSeason(int(latestRace.Season))
	if err != nil {
		log.Println("Unable to load drivers:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}


func SaveDriversController(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/display-leagues", http.StatusSeeOther)
}

func JoinLeagueController(w http.ResponseWriter, r *http.Request) {
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
	leaguePasscode := r.Form.Get("league-passcode")
	if leaguePasscode == "" {
		log.Println("No league passcode provided")
	} else {
		err = service.JoinLeague(user, strings.TrimSpace(leaguePasscode))
		if err != nil {
			log.Println("Error joining league:", err)
		}
	}
	http.Redirect(w, r, "/display-leagues", http.StatusSeeOther)
}

func SaveTeamDetails(w http.ResponseWriter, r *http.Request) {
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
	teamName := r.Form.Get("team-name")
	teamPriciple := r.Form.Get("team-principle")
	err = service.SaveUserTeamDetails(users.User{
		Email: user.Email,
		TeamName: teamName,
		TeamPriciple: teamPriciple,
	})
	if err != nil {
		log.Println("Error saving team details:", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
