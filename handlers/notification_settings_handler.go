package handlers

import (
	"github.com/Simversity/gottp"
	"github.com/parthdesai/sc-notifications/models"
)

type NotificationSettingWithChannelHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingWithChannelHandler) Post(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
	request.Write("Hi!!!!! " + (*channelConfig).UserID)
}

func (self *NotificationSettingWithChannelHandler) Delete(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
}

type NotificationSettingHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingHandler) Post(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
}

func (self *NotificationSettingHandler) Delete(request *gottp.Request) {
	channelConfig := new(models.ChannelConfiguration)
	request.ConvertArguments(channelConfig)
}
