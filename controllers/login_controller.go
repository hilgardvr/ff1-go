package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/users"
	"hilgardvr/ff1-go/view"
	"html/template"
	"log"
	"net/http"
)

func LoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Could not parse form")
	}
	emailAddress := r.Form.Get("email")
	fmt.Println("email: ", emailAddress)
	// templ := "./static/login.html"
	newCode := session.SetLoginCode(emailAddress)
	email.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	// t, err := view.LoginTemplate()//template.ParseFiles(templ)
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	user := users.User{Email: emailAddress, SessionId: ""}
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
	valid := session.ValidateLoginCode(emailAddress, code)
	var templ string
	if valid {
		fmt.Println("successfull code - removing")
		session.DeleteLoginCode(emailAddress)
		session.SetSessionCookie(emailAddress, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		templ = "./static/login.html"
		t, err := template.ParseFiles(templ)
		if err != nil {
			log.Fatalln("template parsing err:", err)
		}
		user := users.User{Email: emailAddress, SessionId: ""}
		err = t.Execute(w, user)
		if err != nil {
			log.Fatalln("template executing err:", err)
		}
		newCode := session.SetLoginCode(emailAddress)
		email.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	}
}