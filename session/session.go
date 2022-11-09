package session

import (
	"errors"
	"fmt"
	"hilgardvr/ff1-go/service"
	"hilgardvr/ff1-go/users"

	"net/http"
	"time"

	"github.com/google/uuid"
)

var expiration = time.Now().Add(1 * time.Minute)

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
	// for _, session := range sessions {
	// 	if session.SessionId == uuid.Value  {
	// 		fmt.Println("Found cookie for ", session.Email)
	// 		return session, nil
	// 	}
	// }
	// err = errors.New("Could not find an existing session")
	// log.Println(err)
	// return users.User{}, err
}

func SetSessionCookie(email string, w http.ResponseWriter) error {
	uuid := uuid.New().String()
	uuidCookie := http.Cookie{Name: "session", Value: uuid, Expires: expiration}
	http.SetCookie(w, &uuidCookie)
	// sessions = append(sessions, users.User{Email: email, SessionId: uuid})
	err := svc.Db.SaveSession(email, uuid)
	fmt.Printf("Cookie set for %s with value %s\n", email, uuid)
	return err
}