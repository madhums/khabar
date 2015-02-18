package main

import (
	"github.com/changer/sc-notifications/config"
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/handlers"
	"gopkg.in/simversity/gottp.v1"
	"log"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	db.Conn = db.GetConn(config.Settings.Sc_Notifications.DBName, config.Settings.Sc_Notifications.DBAddress)
	log.Println("Database Connected :" + config.Settings.Sc_Notifications.DBName + " " + "at address:" + config.Settings.Sc_Notifications.DBAddress)
}

func main() {
	go sysInit()
	registerHandler("notification", "^/notifications/(?P<generic_id>\\w+)/?$", new(handlers.Notification))
	registerHandler("channel", "^/channel/(?P<ident>\\w+)/?$", new(handlers.Gully))
	registerHandler("notification_settings_with_channel", "^/notification_setting/(?P<notification_type>\\w+)/(?P<channel_ident>\\w+)/?$", new(handlers.NotificationSettingWithChannel))
	registerHandler("notification_settings", "^/notification_setting/(?P<type>\\w+)/?$", new(handlers.NotificationSettingHandler))
	registerHandler("User_locale_handler", "^/user_locale/(?P<user>\\w+)/?$", new(handlers.UserLocale))
	gottp.MakeServer(&config.Settings)
}
