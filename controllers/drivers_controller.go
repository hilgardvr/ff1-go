package controllers

import (
	"encoding/json"
	"net/http"
)


func GetDrivers(w http.ResponseWriter, r *http.Request) {
	allDrivers := svc.Db.GetDrivers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDrivers)
}
