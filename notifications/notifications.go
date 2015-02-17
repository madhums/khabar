package notifications

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/parthdesai/sc-notifications/config"
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi/gully"
	"github.com/parthdesai/sc-notifications/dbapi/notification"
	"github.com/parthdesai/sc-notifications/dbapi/notification_instance"
	"github.com/parthdesai/sc-notifications/dbapi/user_locale"
	"log"
	"path"
	"sync"
)

func PrepareTemplateFilename(translationDirectory string, glyIdent string, orgID string, locale string) string {
	return path.Join(translationDirectory, locale+"."+glyIdent+"."+orgID+"."+"json")
}

func PrepareTemplateIdentifier(templateID string, glyIdent string, orgID string, locale string) string {
	return templateID + "." + locale + "." + glyIdent + "." + orgID
}

func SendToAppropriateChannel(dbConn *db.MConn, glyIdent string, user string, appName string, organization string, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	log.Println("Found Channel :" + glyIdent)

	glySetting := gully.FindAppropriateGully(db.DbConnection, user, appName, organization, glyIdent)
	if glySetting == nil {
		log.Println("Unable to find channel")
		return

	}
	userLocale := user_locale.GetFromDatabase(db.DbConnection, user)
	if userLocale == nil {
		log.Println("Unable to find locale for user")
		userLocale = new(user_locale.UserLocale)
		userLocale.Locale = "en_US"
		userLocale.TimeZone = "GMT+0.0"
	}
	filename := PrepareTemplateFilename(config.Settings.Sc_Notifications.TranslationDirectory, glyIdent, organization, userLocale.Locale)
	err := i18n.LoadTranslationFile(filename)
	if err != nil {
		log.Println("Error occured while opening file:" + err.Error())
	}

	T, _ := i18n.Tfunc(userLocale.Locale)

	log.Println(T(PrepareTemplateIdentifier("notification_setting_text", glyIdent, organization, userLocale.Locale), map[string]interface{}{
		"ChannelIdent": glyIdent,
		"AppName":      glySetting.AppName,
		"User":         glySetting.User,
		"Organization": glySetting.Organization,
	}))

}

func SendNotification(dbConn *db.MConn, notificationInstance *notification_instance.NotificationInstance, notificationSetting *notification.Notification) {
	childwg := new(sync.WaitGroup)

	for _, gly := range notificationSetting.Channels {
		go SendToAppropriateChannel(dbConn, gly, notificationInstance.User, notificationInstance.AppName, notificationInstance.Organization, childwg)
	}

	childwg.Wait()
}
