package notifications

import (
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"github.com/changer/sc-notifications/utils"
	"log"
)

var channelMap = map[string]func(*notification_instance.NotificationInstance, string, map[string]interface{}){
	"email": emailChannelHandler,
}

func emailChannelHandler(ntfInst *notification_instance.NotificationInstance, ntfText string, glyData map[string]interface{}) {
	log.Println("Sending email...")

	if ntfInst.Context["Email"] == nil {
		log.Println("Email field not found.")
		return
	}

	email, ok := ntfInst.Context["Email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
	}

	mailConn := utils.MailConn{
		HostName:   glyData["smtp_hostname"].(string),
		UserName:   glyData["smtp_username"].(string),
		Password:   glyData["smtp_password"].(string),
		SenderName: "Changer Spyder",
		Port:       glyData["smtp_port"].(string),
		Host:       glyData["smtp_hostname"].(string) + ":" + glyData["smtp_port"].(string),
	}

	mailConn.SendEmail(utils.Message{
		From:    "no-reply@safetychanger.com",
		To:      []string{email},
		Subject: "Message intended for recipient :" + email + " " + "with name :" + ntfInst.Context["FullName"].(string),
		Body:    ntfText,
	})
}
