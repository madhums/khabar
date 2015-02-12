package notifications

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/utils"
	"log"
	"time"
)

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
	log.Println("Sending notification to " + chanelIdent + " using channel setting _id:" + channelSetting.Id.Hex() + " " + "user id:" + channelSetting.UserID + " " + "app id:" + channelSetting.ApplicationID + " " + "org id:" + channelSetting.OrganizationID)

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
