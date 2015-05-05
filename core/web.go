package core

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/pending"
	"github.com/bulletind/khabar/dbapi/sent"
)

func webHandler(item *pending.PendingItem, text string,
	settings map[string]interface{}) {

	sent_item := db.SentItem{
		CreatedBy:      item.CreatedBy,
		AppName:        item.AppName,
		Organization:   item.Organization,
		User:           item.User,
		IsRead:         false,
		Topic:          item.Topic,
		DestinationUri: item.DestinationUri,
		Text:           text,
		Context:        item.Context,
		Entity:         item.Entity,
	}

	sent_item.PrepareSave()
	sent.Insert(&sent_item)
}
