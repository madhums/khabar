package main

import (
	"github.com/changer/sc-notifications/config"
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/handlers"
	"github.com/changer/sc-notifications/worker"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/simversity/gottp.v2"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	db.Conn = db.GetConn(config.Settings.Sc_Notifications.DBName, config.Settings.Sc_Notifications.DBAddress)
	log.Println("Database Connected :" + config.Settings.Sc_Notifications.DBName + " " + "at address:" + config.Settings.Sc_Notifications.DBAddress)

	transDir := config.Settings.Sc_Notifications.TranslationDirectory

	if len(transDir) == 0 {
		transDir = os.Getenv("PWD") + "/translations"
	}

	log.Println("Directory for translation :" + transDir)

	filepath.Walk(transDir, func(path string, _ os.FileInfo, err error) error {
		fileExt := filepath.Ext(path)
		if fileExt == ".json" && err == nil {
			log.Println("Loading translation file:" + path)
			i18n.LoadTranslationFile(path)
		} else {
			log.Print("Skipping translation file:" + path + " " + "File Extension:" + fileExt + " ")
			if err != nil {
				log.Print("Error:" + err.Error())
			}
		}
		return nil
	})

	log.Println("Translation has been parsed.")

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
