package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"sync"

	"text/template"

	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/bulletind/khabar.v1/config"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/gully"
	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/dbapi/user_locale"
)

const webIdent = "web"
const DEFAULT_LOCALE = "en_US"
const DEFAULT_TIMEZONE = "GMT+0.0"

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

func send(locale, channelIdent string, pending_item *db.PendingItem) {
	if !topics.ChannelAllowed(
		pending_item.User, pending_item.Organization, pending_item.Topic,
		channelIdent,
	) {
		log.Println("Channel", channelIdent, "is blocked for topic", pending_item.Topic)
		return
	}

	var channelData map[string]interface{}

	if channelIdent != WEB {
		channel, err := gully.FindOne(
			pending_item.User,
			pending_item.Organization,
			channelIdent,
		)

		if err != nil {
			log.Println(channelIdent, err.Error())
			return
		}

		channelData = channel.Data
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

func SendNotification(pending_item *db.PendingItem) {
	userLocale, err := user_locale.Get(pending_item.User)
	if err != nil {
		log.Println("Unable to find locale for user", err.Error())
		userLocale = new(db.UserLocale)

		//FIXME:: Please do not hardcode this.
		userLocale.Locale = DEFAULT_LOCALE
		userLocale.TimeZone = DEFAULT_TIMEZONE
	}

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
