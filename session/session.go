package session

import (
	"errors"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/users"

	"net/http"
	"time"

	"github.com/google/uuid"
)

var expiration = time.Hour * 8760
var cookieName = "session"

var svc = service.GetServiceIO()

func GetUserSession(r *http.Request) (users.User, error) {
	uuid, err := r.Cookie(cookieName)
	if err != nil {
		return users.User{}, err
	}
	user, found := svc.Db.GetSession(uuid.Value)
	if !found {
		return users.User{}, errors.New("Session not found")
	}
	latestRace, err := service.GetLatestRace()
	if err != nil {
		return users.User{}, err
	}
	ds, err := svc.Db.GetUserTeamForRace(user, latestRace)
	if err != nil {
		return users.User{}, err
	}
	ls, err := svc.Db.GetLeagueForUser(user)
	user.Team = ds
	user.Leagues = ls
	return user, nil
}

func SetSessionCookie(email string, w http.ResponseWriter) error {
	uuid := uuid.New().String()
	uuidCookie := http.Cookie{Name: cookieName, Value: uuid, Expires: time.Now().Add(expiration)}
	http.SetCookie(w, &uuidCookie)
	err := svc.Db.SaveSession(email, uuid, expiration)
	return err
}

func DeleteUserSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Expires: time.Now()})
}