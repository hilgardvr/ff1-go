package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PickTeamController(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		fmt.Println("no session found")
		err = view.LoginCodeTemplate(w)
		if err != nil {
			log.Println("template executing err:", err)
		}
	} else {
		fmt.Println("session found for email:", user.Email)
		latestRace, err := service.GetLatestRace()
		if err != nil {
			log.Println("Could not get latest race:", err)
			return 
		}
		allDrivers, err := service.GetAllDriversForSeason(int(latestRace.Season))
		if err != nil {
			log.Println("Could not get all drivers by season:", err)
			return 
		}
		err = view.DriversTemplate(w, user, allDrivers)
		if err != nil {
			log.Println("template executing err:", err)
		}
	}
}


func HomeContoller(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		fmt.Println("no session found")
		err = view.LoginCodeTemplate(w)
		if err != nil {
			log.Println("template executing err:", err)
		}
	} else {
		latestCompleted, err := service.GetLatestCompletedRace()
		if err != nil {
			log.Println("could not get latest race:", err)
			return
		}
		latestUserRacePoints, err := service.GetUserRacePoints(user, latestCompleted)
		if err != nil {
			log.Println("could not get latest race points:", err)
			return
		}
		err = view.HomeTemplate(w, user, latestUserRacePoints)
		if err != nil {
			log.Println("template executing err:", err)
		}
	}
}

func LeagueController(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
		return
	}
	err = view.LeagueTemplate(w, user)
	if err != nil {
		log.Println("League template executing err: ", err)
	}

}

func DislayLeagueController(w http.ResponseWriter, r *http.Request) {
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
	query := r.URL.RawQuery
	qs := strings.Split(query, "=")
	if len(qs) != 2 {
		log.Println("Split query incorrect lenth")
		return
	}
	leagueName, err := url.QueryUnescape(qs[0])
	if err != nil {
		log.Println("Unescape query failed: ", err)
		return
	}
	leaguePasscode, err := url.QueryUnescape(qs[1])
	if err != nil {
		log.Println("Unescape query failed: ", err)
		return
	}
	if leaguePasscode == "" {
		log.Println("No league passcode provided")
	} else {
		leagueUsers, err := service.GetLeagueUsers(leaguePasscode)
		if err != nil {
			log.Println("Error fetching league users:", err)
			return
		}
		err = view.DisplayLeagueTemplate(w, user, leagueName, leagueUsers)
		if err != nil {
			log.Println("Error displaying league: ", err)
		}
	}
}

func DisplayRacePoints(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		fmt.Println("no session found")
		err = view.LoginCodeTemplate(w)
		if err != nil {
			log.Println("template executing err:", err)
		}
	} else {
		latestCompleted, err := service.GetLatestCompletedRace()
		if err != nil {
			log.Println("could not get latest race:", err)
			return
		}
		latestUserRacePoints, err := service.GetUserRacePoints(user, latestCompleted)
		if err != nil {
			log.Println("could not get latest race points:", err)
			return
		}
		err = view.RacePointsTemplate(w, user, latestUserRacePoints)
		if err != nil {
			log.Println("template executing err:", err)
		}
	}
}