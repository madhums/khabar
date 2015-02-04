package handlers

import (
	"gopkg.in/simversity/gottp.v1"
)

type NotificationHandler struct {
	gottp.BaseHandler
}

func (self *NotificationHandler) Get(request *gottp.Request) {
	request.Write("Hi!!! " + (*request.UrlArgs)["user_name"])
}

func (self *NotificationHandler) Delete(request *gottp.Request) {
	request.Write("Ohhhhh!!!! Thanks for invoking me. Nobody does that anymore.")
}

func (self *NotificationHandler) Put(request *gottp.Request) {
	request.Write("Ohhhhh!!!! Thanks for invoking me. Nobody does that anymore.")
}
