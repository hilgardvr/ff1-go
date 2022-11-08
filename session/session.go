package session

import (
	"errors"
	"fmt"
	"hilgardvr/ff1-go/users"
	"log"

	"net/http"
	"time"

	"github.com/google/uuid"
)

var expiration = time.Now().Add(1 * time.Minute)

var sessions = []users.User{}

func GetSession(r *http.Request) (users.User, error) {
	uuid, err := r.Cookie("session")
	if err != nil {
		return users.User{}, err
	}
	for _, session := range sessions {
		if session.SessionId == uuid.Value  {
			fmt.Println("Found cookie for ", session.Email)
			return session, nil
		}
	}
	err = errors.New("Could not find an existing session")
	log.Println(err)
	return users.User{}, err
}

func SetSessionCookie(email string, w http.ResponseWriter) {
	uuid := uuid.New().String()
	uuidCookie := http.Cookie{Name: "session", Value: uuid, Expires: expiration}
	http.SetCookie(w, &uuidCookie)
	sessions = append(sessions, users.User{Email: email, SessionId: uuid})
	fmt.Printf("Cookie set for %s with value %s\n", email, uuid)
	return
}