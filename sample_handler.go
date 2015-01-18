package main

import (
	"github.com/Simversity/gottp"
)

type SampleHandler struct {
	gottp.BaseHandler
}

func (self *SampleHandler) Get(request *gottp.Request) {
	request.Write("Hello!!! " + (*request.UrlArgs)["user_name"])
}

func (self *SampleHandler) Patch(request *gottp.Request) {
	request.Write("Ohhhhh!!!! Thanks for invoking me. Nobody does that anymore.")
}
