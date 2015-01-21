package handlers

import (
	"github.com/Simversity/gottp"
	"github.com/parthdesai/sc-notifications/models"
)

type ChannelHandler struct {
	gottp.BaseHandler
}

func (self *ChannelHandler) Post(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
}

func (self *ChannelHandler) Delete(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
}
