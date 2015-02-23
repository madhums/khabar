package notifications

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi/gully"
	"github.com/changer/sc-notifications/dbapi/notification"
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"github.com/changer/sc-notifications/dbapi/user_locale"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"sync"
)

func PrepareTemplateIdentifier(templateID string, glyIdent string) string {
	return templateID + "_" + glyIdent
}

func SendToAppropriateChannel(dbConn *db.MConn, glyIdent string, user string, appName string, org string, destUri string, context map[string]interface{}, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	log.Println("Found Channel :" + glyIdent)

	glySetting := gully.FindAppropriateGully(db.Conn, user, appName, org, glyIdent)
	if glySetting == nil {
		log.Println("Unable to find channel")
		return

	}
	userLocale := user_locale.Get(db.Conn, user)
	if userLocale == nil {
		log.Println("Unable to find locale for user")
		userLocale = new(user_locale.UserLocale)
		userLocale.Locale = "en-US"
		userLocale.TimeZone = "GMT+0.0"
	}

	T, _ := i18n.Tfunc(userLocale.Locale+"_"+appName+"_"+org, userLocale.Locale+"_"+appName, userLocale.Locale)

	context["ChannelIdent"] = glyIdent
	context["AppName"] = appName
	context["User"] = user
	context["Organization"] = org
	context["DestinationUri"] = destUri

	log.Println(T(PrepareTemplateIdentifier("notification_setting_text", glyIdent), context))

}

func SendNotification(dbConn *db.MConn, notificationInstance *notification_instance.NotificationInstance, notificationSetting *notification.Notification) {
	childwg := new(sync.WaitGroup)

	for _, gly := range notificationSetting.Channels {
		go SendToAppropriateChannel(dbConn, gly, notificationInstance.User, notificationInstance.AppName, notificationInstance.Organization, notificationInstance.DestinationUri, notificationInstance.Context, childwg)
	}

	childwg.Wait()
}
