package controllers

import (
	"hilgardvr/ff1-go/drivers"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/session"
	"hilgardvr/ff1-go/view"
	"log"
	"net/http"
	"strconv"
)

var season = 2022

func CreateRacePoints(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
		return
	}
	if user.IsAdmin {
		drivers, err := service.GetAllDriversForSeason(season)
		if err != nil {
			log.Println("error fetching drivers:", err)
			return
		}
		race, err := service.GetLatestRace()
		view.AdminTemplate(w, user, race, drivers)
	} else {
		log.Println("User not an admin:", user.Email)
		http.Redirect(w, r, "/", 300)
	}
}

func UpdateRaceData(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserSession(r)
	if err != nil {
		log.Println("Failed to get user", err)
		return
	}
	if user.IsAdmin {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln("Could not parse form")
		}
		allDrivers, err := service.GetAllDriversForSeason(season)
		if err != nil {
			log.Println("could not get all drivers for admin: ", err)
			return 
		}
		var driverForPoints []drivers.Driver
		for _, v := range allDrivers {
			points := r.Form.Get(strconv.FormatInt(v.Id, 10))
			p, _ := strconv.Atoi(points)
			driverForPoints = append(driverForPoints, drivers.Driver{
				Id: v.Id,
				Name: v.Name,
				Surname: v.Surname,
				Price: v.Price,
				Points: int64(p),
			})
		}
		err = service.CreateRacePoints(driverForPoints)
		if err != nil {
			log.Println("Failed to save driver points:", err)
			return
		}
		http.Redirect(w, r, "/", 300)
	} else {
		log.Println("User not an admin:", user.Email)
		http.Redirect(w, r, "/", 300)
	}
}