package controllers

import (
	"html/template"
	"log"
	"net/http"
	"hilgardvr/ff1-go/email"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	// allDrivers := repo.GetDrivers()
	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	// json, err := json.Marshal(allDrivers)
	// if err != nil {
	// log.Fatalln("template parsing err:", err)
	// }
	// err = t.Execute(w, json)
	err = t.Execute(w, "")
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}

func LoginContoller(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/login.html")
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}

func LoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	emailAddress := r.Form.Get("email")
	log.Println(emailAddress)
	email.SendEmail(emailAddress, "some-code")
	t, err := template.ParseFiles("./static/login_code.html")
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}
