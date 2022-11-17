package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)

const PickTeam = "/pick-team"
const LoginCode = "/logincode"
const Login = "/login"
const Logout = "/logout"
const Home = "/"


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
			err = view.HomeTemplate(w, user)
		} else {
			err = view.HomeTemplate(w, user)
		}
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
	}
}