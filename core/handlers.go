package core

import (
	"os"

	"github.com/bulletind/khabar/db"
)

const (
	EMAIL = "email"
	SMS   = "sms"
	WEB   = "web"
	PUSH  = "push"
	SNS   = "sns"
)

var ChannelMap = map[string]func(*db.PendingItem, string, string, string){
	WEB:   webHandler,
	EMAIL: emailHandler,
	PUSH:  pushHandler,
}

func IsChannelAvailable(ident string) bool {
	_, allowed := ChannelMap[ident]
	return allowed
}

func pushHandler(item *db.PendingItem, text string, locale string, appName string) {
	if len(os.Getenv("SNS_KEY")) == 0 {
		parseHandler(item, text, locale, appName)
	} else {
		snsHandler(item, text, locale, appName)
	}
}
