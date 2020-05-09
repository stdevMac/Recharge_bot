package sendMail

import (
	"log"
	"net/smtp"
)

const (
	/**
		Gmail SMTP Server
	**/
	SMTPServer = "smtp.gmail.com:587"
)

type Sender struct {
	User     string
	Password string
}

func (sender *Sender) SendMail(body , to string)  {

	msg := "From: " + sender.User + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail( SMTPServer,
		smtp.PlainAuth("", sender.User, sender.Password, "smtp.gmail.com"),
		sender.User, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}