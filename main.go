package main

import (
	"github.com/Simversity/gottp"
	"github.com/parthdesai/sc-notifications/handlers"
	"log"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	log.Println("System is ready to Serve")
}

func main() {
	go sysInit()
	registerHandler("notification", "/notifications/(?P<id>\\w+$)", new(handlers.NotificationHandler))
	registerHandler("channel", "/channel/(?P<id>\\w+$)", new(handlers.ChannelHandler))
	registerHandler("notification settings with channel", "/notification_setting/(?P<notification_id>\\w+)/(?P<channel_id>\\w+)/?$", new(handlers.NotificationSettingWithChannelHandler))
	registerHandler("notification settings", "/notification_setting/(?P<notification_id>\\w+)/?$", new(handlers.NotificationSettingHandler))
	gottp.MakeServer(&settings)
}
