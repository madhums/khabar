package utils

import (
	"github.com/bulletind/khabar/config"
	"log"
	"net/smtp"
	"strings"
)

type MailConn struct {
	HostName   string
	UserName   string
	Password   string
	SenderName string
	Port       string
	Host       string
}

type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func (conn *MailConn) getAuth() smtp.Auth {
	return smtp.PlainAuth("", conn.UserName, conn.Password, conn.HostName)
}

func (conn *MailConn) MessageBytes(message Message) []byte {
	subject := "Subject: "
	subject += message.Subject

	subject = strings.TrimSpace(subject)
	from := strings.TrimSpace("From: " + conn.SenderName + " <" + message.From + ">")
	to := strings.TrimSpace("To: " + strings.Join(message.To, ", "))
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	return []byte(subject + "\n" + from + "\n" + to + "\n" + mime +
		"\n\n" + strings.TrimSpace(message.Body))

}

func (conn *MailConn) SendEmail(message Message) {
	err := smtp.SendMail(conn.Host,
		conn.getAuth(),
		message.From,
		message.To,
		conn.MessageBytes(message))

	if err != nil {
		log.Print(err)
		config.Tracer.Notify()
	}
}
