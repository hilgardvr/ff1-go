package controllers

import (
	"encoding/json"
	"net/http"
)


func GetDrivers(w http.ResponseWriter, r *http.Request) {
	allDrivers := svc.Db.GetDrivers()
	// json, err := json.Marshal(allDrivers)
	// if err != nil {
	// 	log.Fatalln("template parsing err:", err)
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}
