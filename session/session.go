package session

import (
	"errors"
	"fmt"
	"hilgardvr/ff1-go/users"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var sessions = []users.User{}
var expiration = time.Now().Add(10 * time.Minute)
var loginCodes = map[string]string{}

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

func SetLoginCode(email string) string {
	if code, found := loginCodes[email]; found {
		return code
	}
	code := rand.Intn(100000)
	str := strconv.Itoa(code)
	padded := fmt.Sprintf("%05s", str)
	loginCodes[email] = padded
	fmt.Println("logincodes:", loginCodes)
	return padded
}

func DeleteLoginCode(email string) {
	delete(loginCodes, email)
}

func ValidateLoginCode(email string, codeToTest string) bool {
	if code, found := loginCodes[email]; found {
		return code == codeToTest
	}
	return false
}
