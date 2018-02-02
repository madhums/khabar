package gotracer

import (
	"log"
	"net/smtp"
	"strings"
)

type mailConn struct {
	Hostname   string
	Username   string
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

func (conn *mailConn) getAuth() smtp.Auth {
	return smtp.PlainAuth("", conn.Username, conn.Password, conn.Hostname)
}

func (conn *mailConn) MessageBytes(message Message) []byte {
	subject := "Subject: "
	subject += message.Subject

	subject = strings.TrimSpace(subject)
	from := strings.TrimSpace("From: " + conn.SenderName + " <" + message.From + ">")
	to := strings.TrimSpace("To: " + strings.Join(message.To, ", "))
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	return []byte(subject + "\n" + from + "\n" + to + "\n" + mime + "\n\n" + strings.TrimSpace(message.Body))

}

func (conn *mailConn) SendEmail(message Message) {
	err := smtp.SendMail(conn.Host,
		conn.getAuth(),
		message.From,
		message.To,
		conn.MessageBytes(message))

	if err != nil {
		log.Panic(err)
	}
}

func MakeConn(settings *Tracer) *mailConn {
	mailconn := &mailConn{
		settings.EmailHost,
		settings.EmailUsername,
		settings.EmailPassword,
		settings.EmailSender,
		settings.EmailPort,
		settings.EmailHost + ":" + settings.EmailPort,
	}
	return mailconn
}
