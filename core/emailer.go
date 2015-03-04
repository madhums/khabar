package core

import (
	"log"

	"github.com/changer/khabar/dbapi/pending"
	"github.com/changer/khabar/utils"
)

func emailer(item *pending.PendingItem, text string, settings map[string]interface{}) {
	log.Println("Sending email...")

	if item.Context["Email"] == nil {
		log.Println("Email field not found.")
		return
	}

	email, ok := item.Context["Email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
	}

	mailConn := utils.MailConn{
		HostName:   settings["smtp_hostname"].(string),
		UserName:   settings["smtp_username"].(string),
		Password:   settings["smtp_password"].(string),
		SenderName: "Changer Spyder",
		Port:       settings["smtp_port"].(string),
		Host:       settings["smtp_hostname"].(string) + ":" + settings["smtp_port"].(string),
	}

	mailConn.SendEmail(utils.Message{
		From:    "no-reply@safetychanger.com",
		To:      []string{email},
		Subject: "Message intended for recipient :" + email + " " + "with name :" + item.Context["FullName"].(string),
		Body:    text,
	})
}
