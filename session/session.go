package session

import (
	"errors"
	"hilgardvr/ff1-go/users"
	"net/http"
	"time"
	"github.com/google/uuid"
)

var sessions = []users.User{}
var expiration = time.Now().Add(365 * 24 * time.Hour)

func GetSession(r *http.Request) (users.User, error) {
	// cookie, _ := r.Cookie(email)
	// log.Println("Found cookie", cookie)
	for _, cookie := range r.Cookies() {
		for _, session := range sessions {
			if session.SessionId == cookie.Value {
				return session, nil
			}
		}
	}
	return users.User{}, errors.New("Could not find an existing session")
}

func SetSession(email string, w http.ResponseWriter) error {
	uuid := uuid.New()
	cookie := http.Cookie{Name: email, Value: uuid.String(), Expires: expiration}
	http.SetCookie(w, &cookie)
	return nil
}