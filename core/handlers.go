package core

import (
	"github.com/bulletind/khabar/dbapi/pending"
)

const (
	EMAIL = "email"
	SMS   = "sms"
	WEB   = "web"
	PUSH  = "push"
)

var ChannelMap = map[string]func(*pending.PendingItem, string,
	map[string]interface{}){
	EMAIL: emailHandler,
	WEB:   webHandler,
	PUSH:  pushHandler,
}

func IsChannelAvailable(ident string) bool {
	_, allowed := ChannelMap[ident]
	return allowed
}
