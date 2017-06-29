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
	"reflect"
	"strings"

	"github.com/aymerick/douceur/inliner"
	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/bulletind/khabar/utils"
	"github.com/scorredoira/email"
)

type mailSettings struct {
	//CSS  string
	BaseTemplate string
	SMTP         *smtpSettings
}

type smtpSettings struct {
	HostName  string
	UserName  string
	Password  string
	Port      string
	FromEmail string
	FromName  string
}

func (smtp smtpSettings) getAddress() string {
	return smtp.HostName + ":" + smtp.Port
}

// load once, store for reuse
var settingsMail *mailSettings

func loadConfig() {
	if settingsMail != nil {
		return
	}

	settingsMail = &mailSettings{
		BaseTemplate: getContentString("email/base.tmpl"),
		//CSS:  getContentString("email/css.css"),
		SMTP: &smtpSettings{
			HostName:  getSMTPEnv("HostName", true),
			UserName:  getSMTPEnv("UserName", true),
			Password:  getSMTPEnv("Password", true),
			Port:      getSMTPEnv("Port", true),
			FromEmail: getSMTPEnv("From_Email", true),
			FromName:  getSMTPEnv("From_Name", false),
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

	var sender string = settingsMail.SMTP.FromName
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

	emailauth := smtp.PlainAuth("", settingsMail.SMTP.UserName, settingsMail.SMTP.Password, settingsMail.SMTP.HostName)
	message := email.NewHTMLMessage(subject, "Dummy")
	message.From = mail.Address{Name: sender, Address: settingsMail.SMTP.FromEmail}
	message.To = []string{emailAddress}
	// inform the user the expected attachments are not there
	item.Context["khabar_attached"] = attachments(item, message)
	message.Body = makeEmail(item, text, locale)

	// send out the email
	err := email.Send(settingsMail.SMTP.getAddress(), emailauth, message)
	if err != nil {
		log.Println("Error sending mail", err)
	}
	saved_item.Insert(db.SavedEmailCollection, &db.SavedItem{Data: message, Details: *item})
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

		// expose all keys to base template as well so we can use them outside the Content as well
		for key, val := range item.Context {
			templateContext[key] = val
		}

		// 1st combine template with css, language specific texts and topic-mail or topic-text
		combined := parse(settingsMail.BaseTemplate, htmlCopy(templateContext))
		// now parse the context from the message
		parsed := parse(combined, htmlCopy(item.Context))
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

func getSMTPEnv(key string, required bool) string {
	return utils.GetEnv("smtp_"+key, required)
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

// copy struct and HTML all string-entries
func htmlCopy(item interface{}) interface{} {
	kind := reflect.TypeOf(item).Kind()
	original := reflect.ValueOf(item)

	if kind == reflect.Slice {
		clone := []interface{}{}
		for i := 0; i < original.Len(); i += 1 {
			clone = append(clone, htmlCopy(original.Index(i).Interface()))
		}
		return clone
	} else if kind == reflect.Map {
		clone := map[string]interface{}{}
		for key, val := range item.(map[string]interface{}) {
			clone[key] = htmlCopy(val)
		}
		return clone
	} else if kind == reflect.String {
		return template.HTML(fmt.Sprint(item))
	} else {
		return item
	}
}

func attachments(item *db.PendingItem, message *email.Message) bool {
	totalSize := int64(0)
	maxSize := int64(8500000) //8.5mb

	attachments := []db.Attachment{}
	for _, attachment := range item.Attachments {
		if strings.HasPrefix(attachment.Type, "image") ||
			strings.HasPrefix(attachment.Type, "audio") ||
			strings.Contains(attachment.Type, "application/vnd.openxmlformats-officedocument") ||
			attachment.Type == "application/pdf" {

			attachments = append(attachments, attachment)
		}
	}

	for _, attachment := range attachments {
		filename, size, err := utils.DownloadFile(attachment.Url, attachment.Name, attachment.IsPrivate)

		if err == nil {
			if (totalSize + size) > maxSize {
				log.Println("Ignoring attachment as email would grow too big", attachment.Url, attachment.Type)
				message.Attachments = make(map[string]*email.Attachment)
				return false
			} else {
				totalSize = totalSize + size
				err = message.Attach(filename)
			}
		}

		if err != nil {
			log.Println("Error attaching file", attachment.Url, err)
		}
	}
	return true
}
