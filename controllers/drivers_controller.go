package controllers

import (
	"encoding/json"
	"hilgardvr/ff1-go/repo"
	"net/http"
)

func GetDrivers(w http.ResponseWriter, r *http.Request) {
	allDrivers := repo.GetDrivers()
	// json, err := json.Marshal(allDrivers)
	// if err != nil {
	// 	log.Fatalln("template parsing err:", err)
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}
