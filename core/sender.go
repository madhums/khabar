package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"text/template"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/processed"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_LOCALE   = "en_US"
	DEFAULT_TIMEZONE = "GMT+0.0"
)

type Parse struct {
	Name string
	Key  string
}

var (
	parseKeys = []Parse{
		Parse{"APP_ID", "parse_application_id"},
		Parse{"API_KEY", "parse_rest_api_key"},
	}

	emailKeys = []string{
		"smtp_hostname",
		"smtp_username",
		"smtp_password",
		"smtp_port",
		"smtp_from",
	}

	locales       = bson.M{}
	localesLoaded = false
)

func getLocales() bson.M {
	if !localesLoaded {
		loadLocales()
	}
	return locales
}

func loadLocales() {
	locales = bson.M{}
	for _, language := range i18n.LanguageTags() {
		// so we load files with names like 'en_US_email', we get 'en-us-email'
		// so we have to make valid stuff again
		key := language[:2] + "-" + strings.ToUpper(language[3:5])
		_, ok := locales[key]
		if !ok {
			value := language[:2] + "_" + strings.ToUpper(language[3:5])
			locales[key] = value
		}
	}

	//fallback for flemish
	_, ok := locales["nl-BE"]
	if !ok {
		locales["nl-BE"] = "nl_NL"
	}

	log.Println("locales", locales)
	localesLoaded = true
}

func sendToChannel(
	pending_item *db.PendingItem,
	text, channelIdent string,
	context map[string]interface{},
) {
	handlerFunc, ok := ChannelMap[channelIdent]
	if !ok {
		log.Println("No handler for Topic:", pending_item.Topic, "Channel:", channelIdent)
		return
	}

	defer config.Tracer.Notify()
	handlerFunc(pending_item, text, context)
}

func getText(locale, ident, channel string, pending_item *db.PendingItem) string {
	T, _ := i18n.Tfunc(
		locale+"_"+pending_item.Organization+"_"+channel,
		locale+"_"+channel,
	)

	text := T(ident, pending_item.Context)
	if text == ident {
		text = ""
	}

	return text
}

// getCategories fetchs distinct available categories to which we can send notifications
func getCategories() []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	var categories []string

	db.Conn.GetCursor(
		session, db.AvailableTopicCollection, utils.M{},
	).Distinct("app_name", &categories)

	return categories
}

// validCategory checks if the category is valid for sending notification
func validCategory(category string) bool {
	categories := getCategories()
	var found bool
	for _, c := range categories {
		if c == category {
			found = true
			break
		}
	}
	return found
}

// getParseKeys returns map of parse api key and app id
// It gets the values from the enviroment variables
func getParseKeys(category string) utils.M {
	doc := utils.M{}

	// Set the Parse api key and id
	for _, parse := range parseKeys {
		envKey := "PARSE_" + category + "_" + parse.Name
		doc[parse.Key] = os.Getenv(envKey)
		if len(os.Getenv(envKey)) == 0 {
			log.Println(envKey, "is empty. Make sure you set this env variable")
		}
	}
	return doc
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

func send(locale, channelName string, pending_item *db.PendingItem) {
	if !topics.ChannelAllowed(
		pending_item.User,
		pending_item.Organization,
		pending_item.AppName,
		pending_item.Topic,
		channelName,
	) {
		log.Println("Channel", channelName, "is blocked for topic", pending_item.Topic)
		return
	}

	if !validCategory(pending_item.AppName) {
		log.Println("Category", pending_item.AppName, "doesn't exist")
		return
	}

	channelData := map[string]interface{}{}

	if channelName == PUSH {
		channelData = getParseKeys(pending_item.AppName)
	}

	text := getText(locale, pending_item.Topic, channelName, pending_item)
	if text == "" {
		// If Topic == text, do not send the notification. This can happen
		// if the translation fails to find a sensible string in the JSON files
		// OR the translation provided was meaningless. To prevent the users
		// from being annpyed, abort this routine.
		log.Println("No translation for:", channelName, pending_item.Topic)
		return
	}

	subject := getText(locale, pending_item.Topic+"_subject", channelName, pending_item)

	if subject != "" {
		pending_item.Context["subject"] = subject
	} else {
		log.Println("Subject not found.")
	}

	if channelName == EMAIL {
		channelData = getEmailKeys()
		buffer := new(bytes.Buffer)

		transDir := config.Settings.Khabar.TranslationDirectory
		path := transDir + "/" + locale + "_base_email.tmpl"

		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Cannot Load the base email template:", path)
		} else {
			t := template.Must(template.New("email").Parse(string(content)))

			data := struct{ Content string }{text}
			t.Execute(buffer, &data)
			text = buffer.String()
		}
	}

	sendToChannel(pending_item, text, channelName, channelData)
}

func ProcessDefaults(user, org string) {
	if !processed.IsProcessed(db.BLANK, org) {
		topics.Initialize(db.BLANK, org)
		processed.MarkAsProcessed(db.BLANK, org)
	}

	if !processed.IsProcessed(user, org) {
		topics.Initialize(user, org)
		processed.MarkAsProcessed(user, org)
	}
}

func SendNotification(pending_item *db.PendingItem) {
	ProcessDefaults(pending_item.User, pending_item.Organization)

	childwg := new(sync.WaitGroup)

	for channel, _ := range ChannelMap {
		childwg.Add(1)

		go func(
			language, channelIdent string,
			pending_item *db.PendingItem,
		) {
			defer childwg.Done()
			send(language, channelIdent, pending_item)
		}(getLocale(pending_item), channel, pending_item)
	}

	childwg.Wait()
}

func getLocale(pending_item *db.PendingItem) string {
	context, ok := pending_item.Context["locale"].(string)
	if ok {
		locale, found := getLocales()[context].(string)
		if found {
			return locale
		}
	}
	return DEFAULT_LOCALE
}
