package controllers

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)

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
}

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