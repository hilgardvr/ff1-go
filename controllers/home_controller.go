package controllers

import (
	"hilgardvr/ff1-go/repo"
	"html/template"
	"log"
	"net/http"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	allDrivers := repo.GetDrivers()
	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Fatalln("template parsing err:", err)
	}
	err = t.Execute(w, allDrivers)
	if err != nil {
		log.Fatalln("template executing err:", err)
	}
}
