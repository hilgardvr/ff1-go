package view

import (
	"fmt"
	"hilgardvr/ff1-go/users"
	"hilgardvr/ff1-go/drivers"
	"html/template"
	"net/http"
)

var basePath = "./view/templates/"
var getLoginCode = basePath + "get_login_code.html"
var login = basePath + "login.html"
var home = basePath + "home.html"
var driversPath = basePath + "drivers.html"
var base = basePath + "base.html"
var league = basePath + "league.html"
var displayLeague = basePath + "display-league.html"
var adminPage = basePath + "admin_page.html"

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
	t, err := template.ParseFiles(base, home)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "base", user) 
	return err
}

func LeagueTemplate(w http.ResponseWriter, user users.User) error {
	fmt.Println(driversPath)
	t, err := template.ParseFiles(base, league)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "base", user)
	return err
}

func DriversTemplate(w http.ResponseWriter, user users.User) error {
	fmt.Println(driversPath)
	t, err := template.ParseFiles(base, driversPath)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "base", user)
	return err
}

func AdminTemplate(w http.ResponseWriter, user users.User, season int, ds []drivers.Driver) error {
	fmt.Println(adminPage)
	t, err := template.ParseFiles(base, adminPage)
	if err != nil {
		return err
	}
	templData := struct {
		Email string
		IsAdmin bool
		Season int
		Drivers []drivers.Driver
	} {
		Email: user.Email,
		IsAdmin: user.IsAdmin,
		Season: season,
		Drivers: ds,
	}
	err = t.ExecuteTemplate(w, "base", templData)
	return err
}

func DisplayLeagueTemplate(w http.ResponseWriter, u users.User, leagueName string, us []users.User) error {
	fmt.Println(displayLeague)
	t, err := template.ParseFiles(base, displayLeague)
	if err != nil {
		return err
	}
	tempalteData := struct {
		Email string
		Users []users.User
		LeagueName string
	}{
		Email: u.Email,
		Users: us,
		LeagueName: leagueName,
	}
	err = t.ExecuteTemplate(w, "base", tempalteData)
	return err
}