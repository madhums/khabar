package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/notifications"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v1"
	"log"
	"net/http"
)

type NotificationHandler struct {
	gottp.BaseHandler
}

func (self *NotificationHandler) Get(request *gottp.Request) {

	notificationInstance := new(models.NotificationInstance)
	request.ConvertArguments(notificationInstance)

	notificationInstance.UserID = request.GetArgument("generic_id").(string)

	paginator := request.GetPaginator()

	request.Write(notificationInstance.GetAllFromDatabase(db.DbConnection, paginator))
}

func (self *NotificationHandler) Put(request *gottp.Request) {
	notificationInstance := new(models.NotificationInstance)
	objectIdString := request.GetArgument("generic_id").(string)
	if !bson.IsObjectIdHex(objectIdString) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Not a valid id."})
		return
	}
	notificationInstance.Id = bson.ObjectIdHex(objectIdString)
	notificationInstance.MarkAsRead(db.DbConnection)
	request.Write(notificationInstance)
}

func (self *NotificationHandler) Post(request *gottp.Request) {
	notificationInstance := new(models.NotificationInstance)
	request.ConvertArguments(notificationInstance)
	notificationInstance.NotificationType = request.GetArgument("generic_id").(string)
	notificationInstance.IsRead = false

	notificationInstance.PrepareSave()

	if !utils.ValidateAndRaiseError(request, notificationInstance) {
		return
	}

	if !notificationInstance.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, ""})
		return
	}

	notificationSetting := models.Notification{
		ApplicationID:  notificationInstance.ApplicationID,
		OrganizationID: notificationInstance.OrganizationID,
		UserID:         notificationInstance.UserID,
		Type:           notificationInstance.NotificationType,
	}

	if !notificationSetting.FindAppropriateNotification(db.DbConnection) {
		log.Println("Unable to find suitable notification setting.")
	} else {
		go notifications.SendNotification(notificationInstance, &notificationSetting, db.DbConnection)
	}

	notificationInstance.InsertIntoDatabase(db.DbConnection)

}
