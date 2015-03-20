package core

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi/gully"
	"github.com/changer/khabar/dbapi/pending"
	"github.com/changer/khabar/dbapi/sent"
	"github.com/changer/khabar/dbapi/topics"
	"github.com/changer/khabar/dbapi/user_locale"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/simversity/gotracer.v1"
	"log"
	"sync"
)

const DEFAULT_LOCALE = "en-US"
const DEFAULT_TIMEZONE = "GMT+0.0"

func sendToChannel(pending_item *pending.PendingItem, text string, channelIdent string, context map[string]interface{}) {
	handlerFunc, ok := channelMap[channelIdent]
	if !ok {
		log.Println("Error : No channel handler found to send this Notification. Topic:" + pending_item.Topic + " Channel:" + channelIdent)
		return
	}
	defer gotracer.Tracer{Dummy: true}.Notify()
	handlerFunc(pending_item, text, context)
}

func send(dbConn *db.MConn, channelIdent string, pending_item *pending.PendingItem) {
	log.Println("Found Channel :" + channelIdent)

	channel := gully.FindOne(db.Conn, pending_item.User, pending_item.AppName, pending_item.Organization, channelIdent)
	if channel == nil {
		log.Println("Unable to find channel")
		return
	}

	userLocale := user_locale.Get(db.Conn, pending_item.User)
	if userLocale == nil {
		log.Println("Unable to find locale for user")
		userLocale = new(user_locale.UserLocale)

		//FIXME:: Please do not hardcode this.
		userLocale.Locale = DEFAULT_LOCALE
		userLocale.TimeZone = DEFAULT_TIMEZONE
	}

	T, _ := i18n.Tfunc(userLocale.Locale+"_"+pending_item.AppName+"_"+pending_item.Organization+"_"+channel.Ident,
		userLocale.Locale+"_"+pending_item.AppName+"_"+channel.Ident, userLocale.Locale+"_"+channel.Ident)

	text := T(pending_item.Topic, pending_item.Context)

	sendToChannel(pending_item, text, channel.Ident, channel.Data)

	sent_item := sent.SentItem{
		AppName:        pending_item.AppName,
		Organization:   pending_item.Organization,
		User:           pending_item.User,
		IsRead:         pending_item.IsRead,
		Topic:          pending_item.Topic,
		DestinationUri: pending_item.DestinationUri,
		Text:           text,
	}

	sent_item.PrepareSave()

	sent.Insert(dbConn, &sent_item)

	log.Println(text)
}

func SendNotification(dbConn *db.MConn,
	pending_item *pending.PendingItem,
	topic *topics.Topic,
) {
	childwg := new(sync.WaitGroup)

	for _, channel := range topic.Channels {
		go func(dbConn *db.MConn, channelIdent string, pending_item *pending.PendingItem, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			send(dbConn, channelIdent, pending_item)
		}(dbConn, channel, pending_item, childwg)
	}

	childwg.Wait()
}
