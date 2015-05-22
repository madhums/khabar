package core

import "gopkg.in/bulletind/khabar.v1/db"

const (
	EMAIL = "email"
	SMS   = "sms"
	WEB   = "web"
	PUSH  = "push"
)

var ChannelMap = map[string]func(*db.PendingItem, string,
	map[string]interface{}){
	EMAIL: emailHandler,
	WEB:   webHandler,
	PUSH:  pushHandler,
}

func IsChannelAvailable(ident string) bool {
	_, allowed := ChannelMap[ident]
	return allowed
}
