package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/users"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
)

var svc = service.GetServiceIO()

func LoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Could not parse form")
	}
	emailAddress := r.Form.Get("email")
	newCode := svc.Db.SetLoginCode(emailAddress)
	svc.EmailService.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	user := users.User{Email: emailAddress}
	err = view.LoginTemplate(w, user)
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Could not parse form")
	}
	code := r.Form.Get("code")
	fmt.Println("code: ", code)
	emailAddress := r.URL.Query().Get("email")
	fmt.Println("email: ", emailAddress)
	valid := svc.Db.ValidateLoginCode(emailAddress, code)
	if valid {
		fmt.Println("successfull code - removing")
		svc.Db.DeleteLoginCode(emailAddress)
		err = session.SetSessionCookie(emailAddress, w)
		if err != nil {
			log.Println("Could not set session:", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		user := users.User{Email: emailAddress}
		err = view.LoginTemplate(w, user)
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
		newCode := svc.Db.SetLoginCode(emailAddress)
		svc.EmailService.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	}
}