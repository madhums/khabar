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

	hasData := notification.GetFromDatabase(db.DbConnection)

	notification.AddChannelToNotification(channelIdent)

	if !hasData {
		log.Println("Creating new document")
		notification.PrepareSave()
		if !notification.IsValid(models.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user_id, org_id and app_id must be present."})
			return
		}
	}

	if !utils.ValidateAndRaiseError(request, notification) {
		log.Println("Validation Failed")
		return
	}

	var err error
	if hasData {
		err = notification.UpdateChannels(db.DbConnection)
	} else {
		log.Println("Successfull call: Inserting document")
		notification.InsertIntoDatabase(db.DbConnection)
	}

	if err != nil {
		log.Println("Error while inserting document :" + err.Error())
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

}

func (self *NotificationSettingWithChannelHandler) Delete(request *gottp.Request) {
	notification := new(models.Notification)

	channelIdent := request.GetArgument("channel_ident").(string)
	notification.Type = request.GetArgument("notification_type").(string)

	request.ConvertArguments(notification)

	hasData := notification.GetFromDatabase(db.DbConnection)

	if !hasData {
		request.Raise(gottp.HttpError{http.StatusNotFound, "notification setting does not exists." + notification.Type})
		return
	}

	notification.RemoveChannelFromNotification(channelIdent)
	log.Println(notification.Channels)

	var err error

	if len(notification.Channels) == 0 {
		err = notification.DeleteFromDatabase(db.DbConnection)
	} else {
		err = notification.UpdateChannels(db.DbConnection)
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
	if !notification.IsValid(models.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}
	err := notification.DeleteFromDatabase(db.DbConnection)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
