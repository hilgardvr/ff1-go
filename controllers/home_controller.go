package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
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
		err = view.DriversTemplate(w, user)
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
		err = view.HomeTemplate(w, user)
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
