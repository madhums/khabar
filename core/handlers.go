package core

import (
	"github.com/changer/khabar/dbapi/db"
	"github.com/changer/khabar/dbapi/pending"
)

var channelMap = map[string]func(*pending.PendingItem, string, map[string]interface{}){
	db.EMAIL: emailer,
	db.WEB:   web_handler,
}
