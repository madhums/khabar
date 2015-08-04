package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"text/template"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/processed"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/dbapi/user_locale"
	"github.com/bulletind/khabar/utils"
	"github.com/nicksnyder/go-i18n/i18n"
)

const webIdent = "web"
const DEFAULT_LOCALE = "en_US"
const DEFAULT_TIMEZONE = "GMT+0.0"

type Parse struct {
	Name string
	Key  string
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
	var keys = []Parse{
		Parse{"APP_ID", "parse_application_id"},
		Parse{"API_KEY", "parse_rest_api_key"},
	}

	// Set the Parse api key and id
	for _, parse := range keys {
		envKey := "PARSE_" + category + "_" + parse.Name
		doc[parse.Key] = os.Getenv(envKey)
		if len(os.Getenv(envKey)) == 0 {
			log.Println("PARSE_"+category+"_"+parse.Name, "is empty. Make sure you set this env variable")
		}
	}
	return doc
}

func send(locale, channelIdent string, pending_item *db.PendingItem) {
	if !topics.ChannelAllowed(
		pending_item.User,
		pending_item.Organization,
		pending_item.AppName,
		pending_item.Topic,
		channelIdent,
	) {
		log.Println("Channel", channelIdent, "is blocked for topic", pending_item.Topic)
		return
	}

	if !validCategory(pending_item.AppName) {
		log.Println("Category", pending_item.AppName, "doesn't exist")
		return
	}

	var channelData utils.M

	if channelIdent != WEB {
		channelData = getParseKeys(pending_item.AppName)
	} else {
		channelData = map[string]interface{}{}
	}

	text := getText(locale, pending_item.Topic, channelIdent, pending_item)
	if text == "" {
		// If Topic == text, do not send the notification. This can happen
		// if the translation fails to find a sensible string in the JSON files
		// OR the translation provided was meaningless. To prevent the users
		// from being annpyed, abort this routine.
		log.Println("No translation for:", channelIdent, pending_item.Topic)
		return
	}

	subject := getText(locale, pending_item.Topic+"_subject", channelIdent, pending_item)

	if subject != "" {
		pending_item.Context["subject"] = subject
	} else {
		log.Println("Subject not found.")
	}

	if channelIdent == EMAIL {
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

	sendToChannel(pending_item, text, channelIdent, channelData)
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
	userLocale, err := user_locale.Get(pending_item.User)
	if err != nil {
		log.Println("Unable to find locale for user", err.Error())
		userLocale = new(db.UserLocale)

		//FIXME:: Please do not hardcode this.
		userLocale.Locale = DEFAULT_LOCALE
		userLocale.TimeZone = DEFAULT_TIMEZONE
	}

	ProcessDefaults(pending_item.User, pending_item.Organization)

	childwg := new(sync.WaitGroup)

	for channel, _ := range ChannelMap {
		childwg.Add(1)

		go func(
			locale, channelIdent string,
			pending_item *db.PendingItem,
		) {
			defer childwg.Done()
			send(locale, channelIdent, pending_item)
		}(userLocale.Locale, channel, pending_item)
	}

	childwg.Wait()
}
