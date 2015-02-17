package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi"
	"github.com/parthdesai/sc-notifications/dbapi/notification"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"log"
	"net/http"
)

type NotificationSettingWithChannelHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingWithChannelHandler) Post(request *gottp.Request) {
	inputNotification := new(notification.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	inputNotification.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(inputNotification)

	inputNotification.AddChannelToNotification(channelIdent)

	ntfication := notification.GetFromDatabase(db.DbConnection, inputNotification.User, inputNotification.ApplicationID, inputNotification.Organization, inputNotification.Type)

	hasData := true

	if ntfication == nil {
		hasData = false
		log.Println("Creating new document")

		inputNotification.PrepareSave()
		if !inputNotification.IsValid(dbapi.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_id must be present."})
			return
		}

		ntfication = inputNotification
	} else {
		ntfication.AddChannelToNotification(channelIdent)
	}

	if !utils.ValidateAndRaiseError(request, ntfication) {
		log.Println("Validation Failed")
		return
	}

	var err error
	if hasData {
		err = notification.UpdateNotification(db.DbConnection, ntfication)
	} else {
		log.Println("Successfull call: Inserting document")
		notification.InsertIntoDatabase(db.DbConnection, ntfication)
	}

	if err != nil {
		log.Println("Error while inserting document :" + err.Error())
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

func (self *NotificationSettingWithChannelHandler) Delete(request *gottp.Request) {
	ntfication := new(notification.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	ntfication.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(ntfication)

	ntfication = notification.GetFromDatabase(db.DbConnection, ntfication.User, ntfication.ApplicationID, ntfication.Organization, ntfication.Type)

	if ntfication == nil {
		request.Raise(gottp.HttpError{http.StatusNotFound, "notification setting does not exists."})
		return
	}

	ntfication.RemoveChannelFromNotification(channelIdent)
	log.Println(ntfication.Channels)

	var err error

	if len(ntfication.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = notification.DeleteFromDatabase(db.DbConnection, ntfication)
	} else {
		log.Println("Updating...")
		err = notification.UpdateNotification(db.DbConnection, ntfication)
	}

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

type NotificationSettingHandler struct {
	gottp.BaseHandler
}

func (self *NotificationSettingHandler) Delete(request *gottp.Request) {
	ntfication := new(notification.Notification)
	request.ConvertArguments(ntfication)
	if !ntfication.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_id must be present."})
		return
	}
	err := notification.DeleteFromDatabase(db.DbConnection, ntfication)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
