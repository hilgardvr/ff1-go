package view

import (
	"encoding/json"
	"fmt"
	"hilgardvr/ff1-go/config"
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/races"
	"hilgardvr/ff1-go/users"
	"html/template"
	"log"
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
var updateRaceMode = basePath + "update_race_disabled.html"

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
	t, err := template.ParseFiles(base, home)
	if err != nil {
		return err
	}
	var teamPrice int64
	for _, v := range user.Team {
		teamPrice += v.Price
	}
	tmpData := struct {
		Email string
		User users.User
		RacePoints races.RacePoints
		TeamPrice int64
	} {
		Email: user.Email,
		User: user,
		RacePoints: latestRacePoints,
		TeamPrice: teamPrice,
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


func RacePointsTemplate(
	w http.ResponseWriter, 
	user users.User, 
	seasonRaces []races.RacePoints,
	userToDisplay users.User,
) error {
	fmt.Println(home)
	t, err := template.ParseFiles(base, racePointsPage)
	if err != nil {
		return err
	}
	tmpData := struct {
		Email string
		User users.User
		RacePoints []races.RacePoints
		UserDisplayed users.User
	} {
		Email: user.Email,
		User: user,
		RacePoints: seasonRaces,
		UserDisplayed: userToDisplay,
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
	tmplData := struct {
		Email string
		User users.User
	} {
		Email: user.Email,
		User: user,
	}
	err = t.ExecuteTemplate(w, "base", tmplData)
	return err
}

func DriversTemplate(w http.ResponseWriter, user users.User, allDrivers []drivers.Driver) error {
	updateMode := config.ServiceConfig().UpdateMode
	if (updateMode) {
		t, err := template.ParseFiles(base, updateRaceMode)
		if err != nil {
			log.Println("Error parsing file:", err)
			return err
		}
		tmplData := struct {
			Email string
			User users.User
		} {
			Email: user.Email,
			User: user,
		}
		err = t.ExecuteTemplate(w, "base", tmplData)
		if err != nil {
			log.Println("Error exicuting:", err)
			return err
		}
		return err
	}
	fmt.Println(driversPath)
	t, err := template.ParseFiles(base, driversPath)
	if err != nil {
		log.Println("Error parsing file:", driversPath)
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
		log.Println("Error marshalling data:", err)
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
	templateData := struct {
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
	err = t.ExecuteTemplate(w, "base", templateData)
	return err
}