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
	"strings"
	"sync"
)

func PrepareTemplateFilename(translationDirectory string, glyIdent string, orgID string, regionID string, languageID string) string {
	return path.Join(translationDirectory, languageID+"-"+strings.ToUpper(regionID)+"."+glyIdent+"."+orgID+"."+"json")
}

func PrepareTemplateIdentifier(templateID string, glyIdent string, orgID string, regionID string, languageID string) string {
	return templateID + "." + languageID + "-" + strings.ToUpper(regionID) + "." + glyIdent + "." + orgID
}

func SendToAppropriateChannel(glyIdent string, userID string, applicationID string, organizationID string, dbConn *db.MConn, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	log.Println("Found Channel :" + glyIdent)

	glySetting := gully.FindAppropriateGully(db.DbConnection, userID, applicationID, organizationID, glyIdent)
	if glySetting == nil {
		log.Println("Unable to find channel")
		return

	}
	userLocale := user_locale.GetFromDatabase(db.DbConnection, userID)
	if userLocale == nil {
		log.Println("Unable to find locale for user:" + userLocale.UserID)
		userLocale = new(user_locale.UserLocale)
		userLocale.RegionID = "US"
		userLocale.LanguageID = "en"
	}
	filename := PrepareTemplateFilename(config.Settings.Sc_Notifications.TranslationDirectory, glyIdent, organizationID, userLocale.RegionID, userLocale.LanguageID)
	err := i18n.LoadTranslationFile(filename)
	if err != nil {
		log.Println("Error occured while opening file:" + err.Error())
	}

	T, _ := i18n.Tfunc(userLocale.LanguageID + "-" + strings.ToUpper(userLocale.RegionID))

	log.Println(T(PrepareTemplateIdentifier("notification_setting_text", glyIdent, organizationID, userLocale.RegionID, userLocale.LanguageID), map[string]interface{}{
		"ChannelIdent":   glyIdent,
		"ApplicationID":  applicationID,
		"UserID":         userID,
		"OrganizationID": organizationID,
	}))

}

func SendNotification(notificationInstance *notification_instance.NotificationInstance, notificationSetting *notification.Notification, dbConn *db.MConn) {
	childwg := new(sync.WaitGroup)

	for _, gly := range notificationSetting.Channels {
		go SendToAppropriateChannel(gly, notificationInstance.UserID, notificationInstance.ApplicationID, notificationInstance.OrganizationID, dbConn, childwg)
	}

	childwg.Wait()
}
