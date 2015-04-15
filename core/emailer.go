package core

import (
	"log"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/pending"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/bulletind/khabar/utils"
)

func emailHandler(
	item *pending.PendingItem, text string, settings map[string]interface{},
) {
	log.Println("Sending email...")

	if item.Context["email"] == nil {
		log.Println("Email field not found.")
		return
	}

	email, ok := item.Context["email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
		//} else {
		//email = "testing@piyushverma.net"
	}

	var sender string = ""
	var subject string = ""

	if item.Context["sender"] != nil {
		sender, ok = item.Context["sender"].(string)
	}

	if item.Context["subject"] != nil {
		subject, ok = item.Context["subject"].(string)
	}

	mailConn := utils.MailConn{
		HostName:   settings["smtp_hostname"].(string),
		UserName:   settings["smtp_username"].(string),
		Password:   settings["smtp_password"].(string),
		SenderName: sender,
		Port:       settings["smtp_port"].(string),
		Host: settings["smtp_hostname"].(string) + ":" +
			settings["smtp_port"].(string),
	}

	msg := utils.Message{
		From:    settings["smtp_from"].(string),
		To:      []string{email},
		Subject: subject,
		Body:    text,
	}

	mailConn.SendEmail(msg)

	saved_item.Insert("email", &db.SavedItem{Data: msg})
}
