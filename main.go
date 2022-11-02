package main

import (
	"hilgardvr/ff1-go/controllers"
	"hilgardvr/ff1-go/repo"
	"log"
	"net/http"
)

const port = ":3000"

func main() {
	err := repo.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server on port", port)
	http.HandleFunc("/", controllers.HomeContoller)
	http.HandleFunc("/logincode", controllers.LoginCodeHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/api/all_drivers", controllers.GetDrivers)
	log.Fatal(http.ListenAndServe(port, nil))
}
