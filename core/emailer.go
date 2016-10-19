package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"reflect"
	"strings"

	"github.com/aymerick/douceur/inliner"
	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/scorredoira/email"
)

type mailSettings struct {
	//CSS  string
	Base string
	SMTP *smtpSettings
}

type smtpSettings struct {
	HostName  string
	UserName  string
	Password  string
	Port      string
	FromEmail string
	FromName  string
}

func (smtp smtpSettings) GetAddress() string {
	return smtp.HostName + ":" + smtp.Port
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
			HostName:  getEnv("HostName", true),
			UserName:  getEnv("UserName", true),
			Password:  getEnv("Password", true),
			Port:      getEnv("Port", true),
			FromEmail: getEnv("From_Email", true),
			FromName:  getEnv("From_Name", false),
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

	emailAddress, ok := item.Context["email"].(string)
	if !ok {
		log.Println("Email field is of invalid type")
		return
	}

	text = makeEmail(item, text, locale)
	var sender string = settings.SMTP.FromName
	var subject string = ""

	if item.Context["sender"] != nil {
		ctxSender, _ := item.Context["sender"].(string)
		if sender != "" {
			sender = fmt.Sprintf("%v (%v)", sender, ctxSender)
		} else {
			sender = ctxSender
		}
	}

	if item.Context["subject"] != nil {
		subject, ok = item.Context["subject"].(string)
	}

	emailauth := smtp.PlainAuth("", settings.SMTP.UserName, settings.SMTP.Password, settings.SMTP.HostName)
	emailContent := email.NewHTMLMessage(subject, text)
	emailContent.From = mail.Address{Name: sender, Address: settings.SMTP.FromEmail}
	emailContent.To = []string{emailAddress}
	//
	//  files := []string{
	//          "big.jpg",
	//          "small.jpg",
	//  } // change here to your own files
	//
	//  for _, filename := range files {
	//          err := emailContent.Attach(filename)
	//
	//          if err != nil {
	//                  fmt.Println(err)
	//          }
	//  }

	// send out the email
	err := email.Send(settings.SMTP.GetAddress(), emailauth, emailContent)
	if err != nil {
		log.Println("Error sending mail", err)
	}
	saved_item.Insert(db.SavedEmailCollection, &db.SavedItem{Data: emailContent, Details: *item})
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
		parsed := parse(combined, copy(item.Context))
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

func getEnv(key string, required bool) string {
	envKey := strings.ToUpper("smtp_" + key)
	value := os.Getenv(envKey)
	if len(os.Getenv(envKey)) == 0 && required {
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

func copy(item interface{}) interface{} {
	kind := reflect.TypeOf(item).Kind()
	original := reflect.ValueOf(item)

	if kind == reflect.Slice {
		clone := []interface{}{}
		for i := 0; i < original.Len(); i += 1 {
			clone = append(clone, copy(original.Index(i).Interface()))
		}
		return clone
	} else if kind == reflect.Map {
		clone := map[string]interface{}{}
		for key, val := range item.(map[string]interface{}) {
			clone[key] = copy(val)
		}
		return clone
	} else if kind == reflect.String {
		return template.HTML(fmt.Sprint(item))
	} else {
		return item
	}
}
