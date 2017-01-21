package core

import (
	"log"
	"strings"
	"sync"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/processed"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"github.com/jinzhu/copier"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_LOCALE   = "en_US"
	DEFAULT_TIMEZONE = "GMT+0.0"
)

var (
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

	localesLoaded = true
}

func sendToChannel(
	pending_item *db.PendingItem,
	text,
	locale,
	appName,
	channelIdent string,
) {
	handlerFunc, ok := ChannelMap[channelIdent]
	if !ok {
		log.Println("No handler for Topic:", pending_item.Topic, "Channel:", channelIdent)
		return
	}

	defer config.Tracer.Notify()
	handlerFunc(pending_item, text, locale, appName)
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

	text := getText(locale, pending_item.Topic, channelName, pending_item)
	if text == "" && channelName != EMAIL {
		// If Topic == text, do not send the notification. This can happen
		// if the translation fails to find a sensible string in the JSON files
		// OR the translation provided was meaningless. To prevent the users
		// from being annoyed, abort this routine.
		log.Println("No translation for:", channelName, pending_item.Topic)
		return
	}

	subject := getText(locale, pending_item.Topic+"_subject", channelName, pending_item)

	if subject != "" {
		pending_item.Context["subject"] = subject
	} else {
		log.Println("Subject not found.")
	}

	linktext := getText(locale, pending_item.Topic+"_linktext", channelName, pending_item)
	if linktext != "" {
		pending_item.Context["linktext"] = linktext
	}

	sendToChannel(pending_item, text, locale, pending_item.AppName, channelName)
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

	locale := getLocale(pending_item)

	childwg := new(sync.WaitGroup)

	user, _ := pending_item.Context["User"].(string)
	fullname, _ := pending_item.Context["fullname"].(string)
	mail, _ := pending_item.Context["email"].(string)
	log.Printf("Sending %s to id:[%s], name:%s, email:%s \n", pending_item.Topic, user, fullname, mail)

	for channel, _ := range ChannelMap {
		childwg.Add(1)

		go func(
			language, channelIdent string,
			pending_item *db.PendingItem,
		) {
			defer childwg.Done()
			// create copy so parallel processes can alter the item
			var clone = new(db.PendingItem)
			copier.Copy(clone, pending_item)

			// maps are not deep cloned - that was what we're after in the end - duh!
			clone.Context = make(map[string]interface{})
			for key, val := range pending_item.Context {
				clone.Context[key] = val
			}

			send(language, channelIdent, clone)
		}(locale, channel, pending_item)
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
