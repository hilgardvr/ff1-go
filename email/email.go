package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var from = os.Getenv("EMAIL_FROM")
var appPassword = os.Getenv("APP_PASSWORD")
var smtpHost = "smtp.gmail.com"
var smtpPort = "587"

func SendEmail(to string, subject string, msg string) error {
	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
	recipients := []string{to}
	message := fmt.Sprintf("Subject: %s\n", subject) + msg
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, []byte(message))
	if err != nil {
		log.Println("failed to send email:", err)
	}
	return err
}
