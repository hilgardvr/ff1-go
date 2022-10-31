package controllers

import (
	"fmt"
	"hilgardvr/ff1-go/email"
	"hilgardvr/ff1-go/session"
	"html/template"
	"log"
	"net/http"
)

//func HomeController(w http.ResponseWriter, r *http.Request) {
//	// allDrivers := repo.GetDrivers()
//	t, err := template.ParseFiles("./static/index.html")
//	if err != nil {
//		log.Fatalln("template parsing err:", err)
//	}
//	// json, err := json.Marshal(allDrivers)
//	// if err != nil {
//	// log.Fatalln("template parsing err:", err)
//	// }
//	// err = t.Execute(w, json)
//	err = t.Execute(w, "")
//	if err != nil {
//		log.Fatalln("template executing err:", err)
//	}
//}

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
	code := r.Form.Get("code")
	fmt.Println(code)
	emailAddress := r.Form.Get("email")
	fmt.Println(emailAddress)
	valid := ValidateUser(code, emailAddress) 
	var templ string
	if (valid) {
		templ = "./static/drivers.html"
		session.SetSession(emailAddress, w)
	} else {
		templ = "./static/login.html"
		newCode := session.SetLoginCode(emailAddress)
		email.SendEmail(emailAddress, newCode)
	}
	t, err := template.ParseFiles(templ)
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	err = t.Execute(w, emailAddress)
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
