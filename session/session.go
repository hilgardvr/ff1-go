package session

import (
	"errors"
	"fmt"
	"hilgardvr/ff1-go/users"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var sessions = []users.User{}
var expiration = time.Now().Add(365 * 24 * time.Hour)
var loginCodes = map[string]string{}

func GetSession(r *http.Request) (users.User, error) {
	// cookie, _ := r.Cookie(email)
	// log.Println("Found cookie", cookie)
	for _, cookie := range r.Cookies() {
		for _, session := range sessions {
			if session.SessionId == cookie.Value {
				fmt.Println("session found:", session)
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

func SetLoginCode(email string) string {
	if code, found := loginCodes[email]; found {
		return code
	}
	code := rand.Int() % 100000
	str := strconv.Itoa(code)
	padded := fmt.Sprintf("%05s", str)
	loginCodes[email] = padded
	fmt.Println("logincodes:", loginCodes)
	return padded
}

func ValidateLoginCode(email string, codeToTest string) bool {
	if code, found := loginCodes[email]; found {
		return code == codeToTest
	}
	return false
}