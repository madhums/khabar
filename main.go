package main

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/handlers"
	"gopkg.in/simversity/gottp.v1"
	"log"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	db.DbConnection = db.GetConn(settings.Sc_Notifications.DBName, settings.Sc_Notifications.DBAddress)
	log.Println("Database Connected :" + settings.Sc_Notifications.DBName + " " + "at address:" + settings.Sc_Notifications.DBAddress)
}

func main() {
	go sysInit()
	registerHandler("notification", "^/notifications/(?P<generic_id>\\w+)/?$", new(handlers.NotificationHandler))
	registerHandler("channel", "^/channel/(?P<ident>\\w+)/?$", new(handlers.ChannelHandler))
	registerHandler("notification settings with channel", "^/notification_setting/(?P<notification_type>\\w+)/(?P<channel_ident>\\w+)/?$", new(handlers.NotificationSettingWithChannelHandler))
	registerHandler("notification settings", "^/notification_setting/(?P<notification_ident>\\w+)/?$", new(handlers.NotificationSettingHandler))
	registerHandler("User locale handler", "^/user_locale/(?P<user_id>\\w+)/?$", new(handlers.UserLocalHandler))
	gottp.MakeServer(&settings)
}
