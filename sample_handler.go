package main

import (
	"github.com/Simversity/gottp"
)

type SampleHandler struct {
	gottp.BaseHandler
}

func (self *SampleHandler) Get(request *gottp.Request) {
	request.Write("Hello!!! \n This is get Request")
}

func (self *SampleHandler) Patch(request *gottp.Request) {
	request.Write("Ohhhhh!!!! Thanks for invoking me. Nobody does that anymore.")
}
