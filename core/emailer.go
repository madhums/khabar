package core

import (
	"github.com/changer/khabar/dbapi/pending"
	"github.com/changer/khabar/utils"
	"log"
)

func emailer(item *pending.PendingItem, text string, settings map[string]interface{}) {
	log.Println("Sending email...")

	if item.Context["email"] == nil {
		log.Println("Email field not found.")
		return
	}

	email, ok := item.Context["email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
	}

	var fullname string = ""
	var sender string = ""
	if item.Context["fullname"] != nil {
		fullname, ok = item.Context["fullname"].(string)
	}

	if item.Context["sender"] != nil {
		sender, ok = item.Context["sender"].(string)
	}

	mailConn := utils.MailConn{
		HostName:   settings["smtp_hostname"].(string),
		UserName:   settings["smtp_username"].(string),
		Password:   settings["smtp_password"].(string),
		SenderName: sender,
		Port:       settings["smtp_port"].(string),
		Host:       settings["smtp_hostname"].(string) + ":" + settings["smtp_port"].(string),
	}

	mailConn.SendEmail(utils.Message{
		From:    settings["smtp_from"].(string),
		To:      []string{email},
		Subject: "Message intended for recipient :" + email + " " + "with name :" + fullname,
		Body:    text,
	})
}
