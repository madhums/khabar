package core

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/bulletind/khabar/utils"
)

var emailKeys = []string{
	"smtp_hostname",
	"smtp_username",
	"smtp_password",
	"smtp_port",
	"smtp_from",
}

func emailHandler(item *db.PendingItem, text string, locale string, appName string) {
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

	settings := getEmailKeys()
	text = makeEmail(text, locale)
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

	saved_item.Insert(db.SavedEmailCollection, &db.SavedItem{Data: msg, Details: *item})
}

func makeEmail(input string, locale string) (output string) {
	buffer := new(bytes.Buffer)

	transDir := config.Settings.Khabar.TranslationDirectory
	path := transDir + "/" + locale + "_base_email.tmpl"

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Cannot Load the base email template:", path)
	} else {
		t := template.Must(template.New("email").Parse(string(content)))

		data := struct{ Content string }{input}
		t.Execute(buffer, &data)
		output = buffer.String()
	}
	return
}

// getEmailKeys returns map of email smtp settings
// It gets the values from the environment variables
func getEmailKeys() utils.M {
	doc := utils.M{}

	// Set the Email key
	for _, key := range emailKeys {
		envKey := strings.ToUpper(key)
		doc[key] = os.Getenv(envKey)
		if len(os.Getenv(envKey)) == 0 {
			log.Println(envKey, "is empty. Make sure you set this env variable")
		}
	}
	return doc
}
