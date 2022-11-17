package view

import (
	"hilgardvr/ff1-go/users"
	"html/template"
	"net/http"
	"fmt"
)

var basePath = "./view/templates/"
var getLoginCode = basePath + "get_login_code.html"
var login = basePath + "login.html"
var home = basePath + "home.html"
var drivers = basePath + "drivers.html"
var base = basePath + "base.html"

func LoginCodeTemplate(w http.ResponseWriter) error {
	fmt.Println(getLoginCode)
	t, err := template.ParseFiles(getLoginCode)
	if err != nil {
		return err
	}
	err = t.Execute(w, "")
	return err
}

func LoginTemplate(w http.ResponseWriter, user users.User) error {
	fmt.Println(login)
	t, err := template.ParseFiles(login)
	if err != nil {
		return err
	}
	err = t.Execute(w, user)
	return err
}

func HomeTemplate(w http.ResponseWriter, user users.User) error {
	fmt.Println(home)
	t, err := template.ParseFiles(home, base)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "base", user) 
	return err
}

func DriversTemplate(w http.ResponseWriter) error {
	fmt.Println(drivers)
	t, err := template.ParseFiles(drivers)
	if err != nil {
		return err
	}
	err = t.Execute(w, "")
	return err
}