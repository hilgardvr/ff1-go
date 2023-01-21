package view

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/races"
	"hilgardvr/ff1-go/users"
	"html/template"
	"net/http"
	"sort"
)

var basePath = "./view/templates/"
var getLoginCode = basePath + "get_login_code.html"
var login = basePath + "login.html"
var home = basePath + "home.html"
var driversPath = basePath + "drivers.html"
var base = basePath + "base.html"
var league = basePath + "league.html"
var displayLeague = basePath + "display-league.html"
var displayLeagues = basePath + "display-leagues.html"
var adminPage = basePath + "admin_page.html"
var racePointsPage = basePath + "race-points.html"

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


func HomeTemplate(w http.ResponseWriter, user users.User, latestRacePoints races.RacePoints) error {
	fmt.Println(home)
	t, err := template.ParseFiles(base, home)
	if err != nil {
		return err
	}
	tmpData := struct {
		Email string
		User users.User
		RacePoints races.RacePoints
	} {
		Email: user.Email,
		User: user,
		RacePoints: latestRacePoints,
	}
	err = t.ExecuteTemplate(w, "base", tmpData) 
	return err
}

func DisplayLeagues(w http.ResponseWriter, user users.User, latestRacePoints races.RacePoints) error {
	fmt.Println(home)
	t, err := template.ParseFiles(base, displayLeagues)
	if err != nil {
		return err
	}
	tmpData := struct {
		Email string
		User users.User
		RacePoints races.RacePoints
	} {
		Email: user.Email,
		User: user,
		RacePoints: latestRacePoints,
	}
	err = t.ExecuteTemplate(w, "base", tmpData) 
	return err
}


func RacePointsTemplate(w http.ResponseWriter, user users.User, seasonRaces []races.RacePoints) error {
	fmt.Println(home)
	t, err := template.ParseFiles(base, racePointsPage)
	if err != nil {
		return err
	}
	tmpData := struct {
		Email string
		User users.User
		RacePoints []races.RacePoints
	} {
		Email: user.Email,
		User: user,
		RacePoints: seasonRaces,
	}
	err = t.ExecuteTemplate(w, "base", tmpData) 
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

func DriversTemplate(w http.ResponseWriter, user users.User, allDrivers []drivers.Driver) error {
	fmt.Println(driversPath)
	t, err := template.ParseFiles(base, driversPath)
	if err != nil {
		return err
	}
	var filteredDrivers []drivers.Driver
	for _, ad := range allDrivers {
		found := false
		for _, ud := range user.Team {
			if ud.Id == ad.Id {
				found = true
			}
		}
		if !found {
			filteredDrivers = append(filteredDrivers, ad)
		}
	}
	userTeam := []drivers.Driver{}
	if user.Team != nil {
		userTeam = user.Team
	}
	data := struct {
		Team []drivers.Driver
		Budget int64
		AllDrivers []drivers.Driver
	} {
		Team: userTeam,
		Budget: 1000000,
		AllDrivers: filteredDrivers,
	}
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	templData := struct {
		Email string 
		User users.User
		TemplData string
	} {
		Email: user.Email,
		User: user,
		TemplData: string(json),
	}
	err = t.ExecuteTemplate(w, "base", templData)
	
	return err
}

func AdminTemplate(w http.ResponseWriter, user users.User, race races.Race, ds []drivers.Driver) error {
	fmt.Println(adminPage)
	t, err := template.ParseFiles(base, adminPage)
	if err != nil {
		return err
	}
	templData := struct {
		Email string
		User users.User
		IsAdmin bool
		Season int64
		Drivers []drivers.Driver
		Race int64
	} {
		Email: user.Email,
		User: user,
		IsAdmin: user.IsAdmin,
		Season: race.Season,
		Drivers: ds,
		Race: race.Race,
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
	sort.Slice(us, func(i int, j int) bool {
		return us[i].SeasonPoints > us[j].SeasonPoints
	})
	tempalteData := struct {
		Email string
		User users.User
		Users []users.User
		LeagueName string
	}{
		Email: u.Email,
		User: u,
		Users: us,
		LeagueName: leagueName,
	}
	err = t.ExecuteTemplate(w, "base", tempalteData)
	return err
}