package core

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/sent"
)

func webHandler(item *db.PendingItem, text string,
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
