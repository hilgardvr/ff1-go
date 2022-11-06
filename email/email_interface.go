package email

import "hilgardvr/ff1-go/config"

type EmailService interface {
	Init(config *config.Config) error
	SendEmail(to string, subject string, msg string) error
}