package email

import (
	"fmt"
	"hilgardvr/ff1-go/config"
	"log"
	"net/smtp"
)

type GmailEmailService struct {
	from string 
	appPassword string
	smtpHost string
	smtpPort string
}

func (g *GmailEmailService) Init(config *config.Config) error {
	// g.from = os.Getenv("EMAIL_FROM")
	// g.appPassword = os.Getenv("APP_PASSWORD")
	g.from = config.EmailFrom
	g.appPassword = config.EmailPassword
	g.smtpHost = "smtp.gmail.com"
	g.smtpPort = "587"
	if g.appPassword == "" || g.from == "" {
		return fmt.Errorf("Could not init gmail from address or password")
	}
	return nil
}

func (g GmailEmailService) SendEmail(to string, subject string, msg string) error {
	auth := smtp.PlainAuth("", g.from, g.appPassword, g.smtpHost)
	recipients := []string{to}
	message := fmt.Sprintf("Subject: %s\n", subject) + msg
	err := smtp.SendMail(g.smtpHost+":"+g.smtpPort, auth, g.from, recipients, []byte(message))
	if err != nil {
		log.Println("failed to send email:", err)
	}
	return err
}
