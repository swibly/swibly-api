package sender

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

type smtpSender struct {
	host     string
	port     string
	email    string
	username string
	password string
}

func (s *smtpSender) Send(to string, subject string, body string) {
	go func() {
		log.Print("INFO: Trying to send an email to `" + to + "` with the subject `" + subject + "`")

		from := mail.Address{Name: s.username, Address: s.email}
		recipient := mail.Address{Address: to}

		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = recipient.String()
		headers["Subject"] = subject

		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + body

		servername := s.host + ":" + s.port
		host, _, err := net.SplitHostPort(servername)
		if err != nil {
			log.Print("ERROR: ", err)
			return
		}

		auth := smtp.PlainAuth("", s.email, s.password, host)

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			log.Print("ERROR: ", err)
			return
		}

		c, err := smtp.NewClient(conn, host)
		defer c.Close()

		if err != nil {
			log.Print("ERROR: ", err)
			return
		}

		if err = c.Auth(auth); err != nil {
			log.Print("ERROR: ", err)
			return
		}

		if err = c.Mail(from.Address); err != nil {
			log.Print("ERROR: ", err)
			return
		}

		if err = c.Rcpt(recipient.Address); err != nil {
			log.Print("ERROR: ", err)
			return
		}

		w, err := c.Data()
		defer w.Close()

		if err != nil {
			log.Print("ERROR: ", err)
			return
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			log.Print("ERROR: ", err)
			return
		}

		log.Print("INFO: An email has been sent to `" + to + "` with the subject `" + subject + "`")
	}()
}
