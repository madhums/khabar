package notifications

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi/gully"
	"github.com/changer/sc-notifications/dbapi/notification"
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"github.com/changer/sc-notifications/dbapi/sent_notification"
	"github.com/changer/sc-notifications/dbapi/user_locale"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"sync"
)

func PrepareTemplateIdentifier(templateID string, glyIdent string) string {
	return templateID + "_" + glyIdent
}

func SendToAppropriateChannel(dbConn *db.MConn, glyIdent string, ntfInst *notification_instance.NotificationInstance, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	log.Println("Found Channel :" + glyIdent)

	glySetting := gully.FindAppropriateGully(db.Conn, ntfInst.User, ntfInst.AppName, ntfInst.Organization, glyIdent)
	if glySetting == nil {
		log.Println("Unable to find channel")
		return

	}
	userLocale := user_locale.Get(db.Conn, ntfInst.User)
	if userLocale == nil {
		log.Println("Unable to find locale for user")
		userLocale = new(user_locale.UserLocale)
		userLocale.Locale = "en-US"
		userLocale.TimeZone = "GMT+0.0"
	}

	T, _ := i18n.Tfunc(userLocale.Locale+"_"+ntfInst.AppName+"_"+ntfInst.Organization, userLocale.Locale+"_"+ntfInst.AppName, userLocale.Locale)

	ntfInst.Context["ChannelIdent"] = glyIdent
	ntfInst.Context["AppName"] = ntfInst.AppName
	ntfInst.Context["User"] = ntfInst.User
	ntfInst.Context["Organization"] = ntfInst.Organization
	ntfInst.Context["DestinationUri"] = ntfInst.DestinationUri

	ntfText := T(PrepareTemplateIdentifier("notification_setting_text", glyIdent), ntfInst.Context)

	sentNtf := sent_notification.NotificationInstance{
		AppName:          ntfInst.AppName,
		Organization:     ntfInst.Organization,
		User:             ntfInst.User,
		IsRead:           ntfInst.IsRead,
		NotificationType: ntfInst.NotificationType,
		DestinationUri:   ntfInst.DestinationUri,
		NotificationText: ntfText,
	}

	sentNtf.PrepareSave()

	sent_notification.Insert(dbConn, &sentNtf)

	log.Println(ntfText)

}

func SendNotification(dbConn *db.MConn, ntfInst *notification_instance.NotificationInstance, ntfSetting *notification.Notification) {
	childwg := new(sync.WaitGroup)

	for _, gly := range ntfSetting.Channels {
		go SendToAppropriateChannel(dbConn, gly, ntfInst, childwg)
	}

	childwg.Wait()
}
