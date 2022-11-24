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

const PickTeam = "/pick-team"
const TeamName = "/submit-team-name"
const LoginCode = "/logincode"
const Login = "/login"
const Logout = "/logout"
const Home = "/"

var svc = service.GetServiceIO()

func SaveController(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
	}
	var drivers []drivers.Driver
	err = json.NewDecoder(r.Body).Decode(&drivers)
	if err != nil {
		log.Println("Failed to parse drivers", err)
	}
	fmt.Println(user)
	fmt.Println(drivers)
}

func HomeContoller(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		fmt.Println("no session found")
		err = view.LoginCodeTemplate(w)
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
	} else {
		fmt.Println("session found for email:", user.Email)
		if r.URL.Path == PickTeam {
			fmt.Println("Picking team")
			err = view.DriversTemplate(w, user)
		} else if r.URL.Path == TeamName {
			svc.Db.SaveTeamName(user, r.Form.Get("team-name"))
			err = view.HomeTemplate(w, user)
		} else {
			err = view.HomeTemplate(w, user)
		}
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
	}
}