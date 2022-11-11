package session

import (
	"errors"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/users"

	"net/http"
	"time"

	"github.com/google/uuid"
)

var expiration = time.Minute * 2

var svc = service.GetServiceIO()

func GetSession(r *http.Request) (users.User, error) {
	uuid, err := r.Cookie("session")
	if err != nil {
		return users.User{}, err
	}
	email, found := svc.Db.GetSession(uuid.Value)
	if !found {
		return users.User{}, errors.New("Session not found")
	}
	return users.User{Email: email}, nil
}

func SetSessionCookie(email string, w http.ResponseWriter) error {
	uuid := uuid.New().String()
	uuidCookie := http.Cookie{Name: "session", Value: uuid, Expires: time.Now().Add(expiration)}
	http.SetCookie(w, &uuidCookie)
	err := svc.Db.SaveSession(email, uuid, expiration)
	return err
}