package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)


func HomeContoller(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetSession(r)
	if err != nil {
		fmt.Println("no session found")
		err = view.LoginCodeTemplate(w)
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
	} else {
		fmt.Println("session found for email:", user.Email)
		err = view.HomeTemplate(w, user)
	}
}