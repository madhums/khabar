package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aymerick/douceur/inliner"
	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/bulletind/khabar/utils"
)

type mailSettings struct {
	//CSS  string
	Base string
	SMTP *smtpSettings
}

type smtpSettings struct {
	HostName string
	UserName string
	Password string
	Port     string
	From     string
}

// load once, store for reuse
var settings *mailSettings

func loadConfig() {
	if settings != nil {
		return
	}

	settings = &mailSettings{
		Base: getContentString("email/base.tmpl"),
		//CSS:  getContentString("email/css.css"),
		SMTP: &smtpSettings{
			HostName: getEnv("HostName"),
			UserName: getEnv("UserName"),
			Password: getEnv("Password"),
			Port:     getEnv("Port"),
			From:     getEnv("From"),
		},
	}
}

func emailHandler(item *db.PendingItem, text string, locale string, appName string) {
	log.Println("Sending email...")
	loadConfig()

	if item.Context["email"] == nil {
		log.Println("Email field not found.")
		return
	}

	email, ok := item.Context["email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
	}

	text = makeEmail(item, text, locale)
	var sender string = ""
	var subject string = ""

	if item.Context["sender"] != nil {
		sender, ok = item.Context["sender"].(string)
	}

	if item.Context["subject"] != nil {
		subject, ok = item.Context["subject"].(string)
	}

	mailConn := utils.MailConn{
		HostName:   settings.SMTP.HostName,
		UserName:   settings.SMTP.UserName,
		Password:   settings.SMTP.Password,
		Port:       settings.SMTP.Port,
		Host:       settings.SMTP.HostName + ":" + settings.SMTP.Port,
		SenderName: sender,
	}

	msg := utils.Message{
		From:    settings.SMTP.From,
		To:      []string{email},
		Subject: subject,
		Body:    text,
	}

	mailConn.SendEmail(msg)

	saved_item.Insert(db.SavedEmailCollection, &db.SavedItem{Data: msg, Details: *item})
}

func makeEmail(item *db.PendingItem, topicMail string, locale string) string {
	// get json translations for template
	templateContext := getTemplateContext(locale)

	if topicMail == "" {
		topicMail = getContentString(fmt.Sprintf("%v_email/%v.tmpl", locale, item.Topic))
	}

	if templateContext != nil && topicMail != "" {
		templateContext["Content"] = template.HTML(topicMail)
		//templateContext["CSS"] = settings.CSS

		subject, ok := item.Context["subject"].(string)
		if ok && subject != "" {
			templateContext["Subject"] = subject
		}

		// 1st combine template with css, language specifixc texts and topic-mail or topic-text
		combined := parse(settings.Base, templateContext)
		// now parse the context from the message
		parsed := parse(combined, item.Context)
		// and change from css to style per element
		output, err := inliner.Inline(parsed)
		if err != nil {
			log.Println("Error parsing css:", err)
		}
		return output
	}
	return ""
}

func getTemplateContext(locale string) map[string]interface{} {
	// get json translations for template
	var templateContext map[string]interface{}

	localeContext := getContent(fmt.Sprintf("%v_base_email.json", locale))
	if localeContext == nil {
		log.Println("No locale " + locale + " context found for template:")
		return templateContext
	}

	// parse to json
	err := json.Unmarshal(localeContext, &templateContext)
	if err != nil {
		log.Println("Error parsing locale context to json:", err)
	}
	return templateContext
}

func getEnv(key string) string {
	envKey := strings.ToUpper("smtp_" + key)
	value := os.Getenv(envKey)
	if len(os.Getenv(envKey)) == 0 {
		log.Println(envKey, "is empty. Make sure you set this env variable")
	}
	return value
}

func getContentString(subpath string) string {
	return string(getContent(subpath))
}

func getContent(subpath string) (output []byte) {
	transDir := config.Settings.Khabar.TranslationDirectory
	path := transDir + "/" + subpath
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Cannot Load the template:", path)
	} else {
		output = content
	}
	return
}

func parse(content string, data interface{}) string {
	buffer := new(bytes.Buffer)
	t := template.Must(template.New("email").Parse(string(content)))
	t.Execute(buffer, &data)
	return buffer.String()
}
