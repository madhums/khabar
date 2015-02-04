package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"log"
	"net/http"
)

type NotificationSettingWithChannelHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingWithChannelHandler) Post(request *gottp.Request) {
	notification := new(models.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	notification.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(notification)

	hasData := db.DbConnection.Get("notifications", map[string]interface{}{"app_id": notification.ApplicationID,
		"org_id": notification.OrganizationID, "user_id": notification.UserID, "type": notification.Type}).Next(notification)

	newChannels := make([]string, len(notification.Channels)+1)
	copy(newChannels, notification.Channels)
	newChannels[len(newChannels)-1] = channelIdent
	notification.Channels = newChannels

	utils.RemoveDuplicates(&(notification.Channels))

	if !hasData {
		log.Println("Creating new document")
		notification.PrepareSave(db.DbConnection)
	}

	if !utils.ValidateAndRaiseError(request, notification) {
		return
	}

	var err error
	if hasData {

		err = db.DbConnection.Update("notifications", db.M{"_id": notification.Id},
			db.M{
				"$set": db.M{
					"channels": notification.Channels,
				},
			})

	} else {
		db.DbConnection.Insert("notifications", notification)
	}

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

func (self *NotificationSettingWithChannelHandler) Delete(request *gottp.Request) {
	notification := new(models.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	notification.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(notification)

	hasData := db.DbConnection.Get("notifications", map[string]interface{}{"app_id": notification.ApplicationID,
		"org_id": notification.OrganizationID, "user_id": notification.UserID, "type": notification.Type}).Next(notification)

	if !hasData {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, "notification setting does not exists." + notification.Type})
		return
	}

	utils.RemoveElement(&(notification.Channels), channelIdent)

	var err error

	if len(notification.Channels) == 0 {

		err = db.DbConnection.Delete("notifications", map[string]interface{}{"_id": notification.Id})

	} else {

		err = db.DbConnection.Update("notifications", db.M{"_id": notification.Id},
			db.M{
				"$set": db.M{
					"channels": notification.Channels,
				},
			})
	}

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

type NotificationSettingHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingHandler) Delete(request *gottp.Request) {
	notification := new(models.Notification)
	request.ConvertArguments(notification)
	if !notification.IsValid() {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}
	err := db.DbConnection.Delete("notifications", map[string]interface{}{"app_id": notification.ApplicationID,
		"org_id": notification.OrganizationID, "user_id": notification.UserID, "type": notification.Type})
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
