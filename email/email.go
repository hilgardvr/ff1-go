package email

import (
	"log"
	"net/smtp"
	"os"
)

var from = os.Getenv("EMAIL_FROM")
var appPassword = os.Getenv("APP_PASSWORD")
var smtpHost = "smtp.gmail.com"
var smtpPort = "587"


func SendEmail(to string, msg string) {
	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
	recipients := []string{to}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, []byte(msg))
	if err != nil {
		log.Println("failed to send email:", err)
	}
	return
}