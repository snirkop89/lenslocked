package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@lenslocked.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return &es
}

type EmailService struct {
	// DefaultSender is used as the default sender when oneisn't provided for an
	// email. This is also used in function where the email is predeterminded,
	// like the forgoteen password email.
	DefaultSender string

	// unexported fields
	dialer *mail.Dialer
}
