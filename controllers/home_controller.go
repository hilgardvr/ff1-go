package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/users"
	"html/template"
	"log"
	"net/http"
)

func HomeContoller(w http.ResponseWriter, r *http.Request) {
	session, err := session.GetSession(r)
	var templ string
	if err != nil {
		fmt.Println("no session found")
		templ = "./static/signin.html"
	} else {
		fmt.Println("session found for email:", session.Email)
		templ = "./static/drivers.html"
	}
	t, err := template.ParseFiles(templ)
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Fatalln("template executing err:", err)
	}

}

func LoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Could not parse form")
	}
	emailAddress := r.Form.Get("email")
	fmt.Println("email: ", emailAddress)
	templ := "./static/login.html"
	newCode := session.SetLoginCode(emailAddress)
	email.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	t, err := template.ParseFiles(templ)
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	user := users.User{Email: emailAddress, SessionId: ""}
	err = t.Execute(w, user)
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
	valid := ValidateUser(code, emailAddress)
	var templ string
	if valid {
		fmt.Println("successfull code")
		templ = "./static/drivers.html"
		session.SetSessionCookie(emailAddress, w)
	} else {
		templ = "./static/login.html"
		newCode := session.SetLoginCode(emailAddress)
		email.SendEmail(emailAddress, "Your F1-Go login code", newCode)
	}
	t, err := template.ParseFiles(templ)
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	user := users.User{Email: emailAddress, SessionId: ""}
	err = t.Execute(w, user)
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}

func ValidateUser(code, email string) bool {
	if code == "" || email == "" {
		return false
	}
	return session.ValidateLoginCode(email, code)
}
