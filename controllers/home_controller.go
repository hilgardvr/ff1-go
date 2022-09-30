package controllers

import (
	"html/template"
	"log"
	"net/http"
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
