package controllers

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/users"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)

const PickTeam = "/pick-team"
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
	var ds []drivers.Driver
	err = json.NewDecoder(r.Body).Decode(&ds)
	if err != nil {
		log.Println("Failed to parse drivers", err)
	}
	valid := drivers.ValidateTeam(ds)
	if valid {
		err = svc.Db.SaveTeam(user, ds)
		if err != nil {
			log.Println("Failed to save team for user: ", user, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return
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
		fmt.Println("session found for email:", user.Email)
		if r.URL.Path == PickTeam {
			fmt.Println("Picking team")
			err = view.DriversTemplate(w, user)
		} else {
			team, err := svc.Db.GetTeam(user)
			if err != nil {
				log.Println("Failed to fetch user team for user: ", user, err)
				return
			}
			if len(team) > 0 {
				user = users.User{
					Email: user.Email,
					Team: team,
				}
			}
			err = view.HomeTemplate(w, user)
		}
		if err != nil {
			log.Println("template executing err:", err)
		}
	}
}