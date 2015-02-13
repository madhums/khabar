package notifications

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/utils"
	"log"
	"os"
	"strings"
	"time"
)

func PrepareTemplateIdentifier(channelIdent string, orgID string, regionID string, languageID string) string {
	return languageID + "-" + strings.ToUpper(regionID) + "." + channelIdent + "." + orgID + "." + "json"
}

func SendToAppropriateChannel(chanelIdent string, applicationID string, organizationID string, userID string, dbConn *db.MConn, wg *utils.TimedWaitGroup) {

	wg.Add(1)
	defer wg.Done()

	log.Println("Found Channel :" + chanelIdent)
	channelSetting := models.Channel{
		Ident:          chanelIdent,
		ApplicationID:  applicationID,
		OrganizationID: organizationID,
		UserID:         userID,
	}
	channelSetting.FindAppropriateChannel(dbConn)
	userLocale := models.UserLocale{
		UserID: userID,
	}
	if !userLocale.GetFromDatabase(dbConn) {
		log.Println("Unable to find locale for user:" + userLocale.UserID)
		userLocale.RegionID = "US"
		userLocale.LanguageID = "en"
	}
	filename := PrepareTemplateIdentifier(chanelIdent, organizationID, userLocale.RegionID, userLocale.LanguageID)
	err := i18n.LoadTranslationFile(filename)
	if err != nil {
		log.Println(os.Getenv("PWD"))
		log.Println("Error occured while opening file:" + err.Error())
	}

	T, _ := i18n.Tfunc(userLocale.LanguageID + "-" + strings.ToUpper(userLocale.RegionID))

	log.Println(T("notification_setting_text", map[string]interface{}{
		"ChannelIdent":   chanelIdent,
		"ApplicationID":  applicationID,
		"UserID":         userID,
		"OrganizationID": organizationID,
	}))

}

func SendNotification(notificationInstance *models.NotificationInstance, notificationSetting *models.Notification, dbConn *db.MConn) {
	childwg := new(utils.TimedWaitGroup)
	childwg.TimeOut = 5 * time.Minute

	for _, channel := range notificationSetting.Channels {
		go SendToAppropriateChannel(channel, notificationInstance.ApplicationID, notificationInstance.OrganizationID, notificationInstance.UserID, dbConn, childwg)
	}

	hasCompletedSuccessfully := childwg.TimedWait()
	if !hasCompletedSuccessfully {
		log.Println("Goroutine spanwed to send notification instance id:=" + notificationInstance.Id.Hex() + " " + "and notification settings id:=" + notificationSetting.Id.Hex() + " " + "was timedout.")
	}

}
