package sender

import "github.com/swibly/swibly-api/config"

type sender interface {
	Send(to string, subject string, body string) error
}

var SMTPSender *smtpSender

func Init() {
	SMTPSender = &smtpSender{
		host:     config.SMTP.Host,
		port:     config.SMTP.Port,
		username: config.SMTP.Username,
		email:    config.SMTP.Email,
		password: config.SMTP.Password,
	}
}
