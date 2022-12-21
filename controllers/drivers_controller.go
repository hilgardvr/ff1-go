package controllers

import (
	"encoding/json"
	"hilgardvr/ff1-go/service"
	"log"
	"net/http"
)


func GetDrivers(w http.ResponseWriter, r *http.Request) {
	allDrivers, err := service.GetAllDrivers()
	if err != nil {
		log.Println("Unable to load drivers:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}

func GetBudget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(1000000)
}