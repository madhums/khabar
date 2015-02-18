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

type NotificationSettingWithChannel struct {
	gottp.BaseHandler
}

func (self *NotificationSettingWithChannel) Post(request *gottp.Request) {
	inputNotification := new(notification.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	inputNotification.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(inputNotification)

	inputNotification.AddChannel(channelIdent)

	ntfication := notification.Get(db.Conn, inputNotification.User, inputNotification.AppName, inputNotification.Organization, inputNotification.Type)

	hasData := true

	if ntfication == nil {
		hasData = false
		log.Println("Creating new document")

		inputNotification.PrepareSave()
		if !inputNotification.IsValid(dbapi.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
			return
		}

		ntfication = inputNotification
	} else {
		ntfication.AddChannel(channelIdent)
	}

	if !utils.ValidateAndRaiseError(request, ntfication) {
		log.Println("Validation Failed")
		return
	}

	var err error
	if hasData {
		err = notification.Update(db.Conn, ntfication.User, ntfication.AppName, ntfication.Organization, ntfication.Type, &db.M{
			"channels": ntfication.Channels,
		})
	} else {
		log.Println("Successfull call: Inserting document")
		notification.Insert(db.Conn, ntfication)
	}

	if err != nil {
		log.Println("Error while inserting document :" + err.Error())
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

func (self *NotificationSettingWithChannel) Delete(request *gottp.Request) {
	ntfication := new(notification.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	ntfication.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(ntfication)

	ntfication = notification.Get(db.Conn, ntfication.User, ntfication.AppName, ntfication.Organization, ntfication.Type)

	if ntfication == nil {
		request.Raise(gottp.HttpError{http.StatusNotFound, "notification setting does not exists."})
		return
	}

	ntfication.RemoveChannel(channelIdent)
	log.Println(ntfication.Channels)

	var err error

	if len(ntfication.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = notification.Delete(db.Conn, &db.M{"app_name": ntfication.AppName,
			"org": ntfication.Organization, "user": ntfication.User, "type": ntfication.Type})
	} else {
		log.Println("Updating...")
		err = notification.Update(db.Conn, ntfication.User, ntfication.AppName, ntfication.Organization, ntfication.Type, &db.M{
			"channels": ntfication.Channels,
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
	ntfication := new(notification.Notification)
	request.ConvertArguments(ntfication)
	if !ntfication.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}
	err := notification.Delete(db.Conn, &db.M{"app_name": ntfication.AppName,
		"org": ntfication.Organization, "user": ntfication.User, "type": ntfication.Type})
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
