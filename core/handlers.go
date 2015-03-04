package core

import "github.com/changer/khabar/dbapi/pending"

var channelMap = map[string]func(*pending.PendingItem, string, map[string]interface{}){
	"email": emailer,
}
