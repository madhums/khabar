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
	"github.com/bulletind/khabar/dbapi/gully"
	"github.com/bulletind/khabar/dbapi/pending"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/dbapi/user_locale"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/simversity/gotracer.v1"
)

const webIdent = "web"
const DEFAULT_LOCALE = "en_US"
const DEFAULT_TIMEZONE = "GMT+0.0"

func sendToChannel(
	pending_item *pending.PendingItem,
	text, channelIdent string,
	context map[string]interface{},
) {
	handlerFunc, ok := ChannelMap[channelIdent]
	if !ok {
		log.Println("No handler for Topic:" + pending_item.Topic + " Channel:" + channelIdent)
		return
	}

	defer gotracer.Tracer{Dummy: true}.Notify()
	handlerFunc(pending_item, text, context)
}

func getText(locale, ident, channel string, pending_item *pending.PendingItem) string {
	T, _ := i18n.Tfunc(
		locale+"_"+pending_item.AppName+"_"+pending_item.Organization+"_"+channel,
		locale+"_"+pending_item.AppName+"_"+channel,
		locale+"_"+channel,
	)

	text := T(ident, pending_item.Context)
	if text == ident {
		text = ""
	}

	return text
}

func send(locale, channelIdent string, pending_item *pending.PendingItem) {

	if !topics.ChannelAllowed(pending_item.User, pending_item.AppName,
		pending_item.Organization, pending_item.Topic, channelIdent) {
		log.Println("Channel :" + channelIdent + " " + "is blocked for topic :" + pending_item.Topic)
		return
	}

	channel, err := gully.FindOne(
		pending_item.User,
		pending_item.AppName, pending_item.Organization,
		channelIdent,
	)

	if err != nil {
		log.Println("Unable to find channel : " + channelIdent + err.Error())
		return
	}

	text := getText(locale, pending_item.Topic, channelIdent, pending_item)
	if text == "" {
		// If Topic == text, do not send the notification. This can happen
		// if the translation fails to find a sensible string in the JSON files
		// OR the translation provided was meaningless. To prevent the users
		// from being annpyed, abort this routine.

		log.Println("No translation for:" + channelIdent + pending_item.Topic)
		return
	}

	if channelIdent == EMAIL || channel.Ident == PUSH {
		subject := getText(locale, pending_item.Topic+"_subject", channelIdent, pending_item)

		if subject != "" {
			pending_item.Context["subject"] = subject
		}
	}

	if channelIdent == EMAIL {
		buffer := new(bytes.Buffer)

		transDir := config.Settings.Khabar.TranslationDirectory
		path := transDir + "/" + locale + "_base_email.tmpl"

		if _, err := os.Stat(path); err == nil {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Cannot Load the base email template")
			} else {
				t := template.Must(template.New("email").Parse(string(content)))

				data := struct{ Content string }{text}
				t.Execute(buffer, &data)
				text = buffer.String()
			}
		}
	}

	sendToChannel(pending_item, text, channel.Ident, channel.Data)
}

func SendNotification(pending_item *pending.PendingItem) {
	userLocale, err := user_locale.Get(pending_item.User)
	if err != nil {
		log.Println("Unable to find locale for user :" + err.Error())
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
			pending_item *pending.PendingItem,
		) {
			defer childwg.Done()
			send(locale, channelIdent, pending_item)
		}(userLocale.Locale, channel, pending_item)
	}

	childwg.Wait()
}
