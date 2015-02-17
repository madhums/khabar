package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi/notification"
	"github.com/parthdesai/sc-notifications/dbapi/notification_instance"
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

	notificationInstance := new(notification_instance.NotificationInstance)
	request.ConvertArguments(notificationInstance)

	notificationInstance.User = request.GetArgument("generic_id").(string)

	paginator := request.GetPaginator()

	request.Write(notification_instance.GetAllFromDatabase(db.DbConnection, paginator, notificationInstance.User, notificationInstance.ApplicationID, notificationInstance.Organization))
}

func (self *NotificationHandler) Put(request *gottp.Request) {
	notificationInstance := new(notification_instance.NotificationInstance)
	objectIdString := request.GetArgument("generic_id").(string)
	if !bson.IsObjectIdHex(objectIdString) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "_id is not a valid Hex object."})
		return
	}
	notificationInstance.Id = bson.ObjectIdHex(objectIdString)
	notification_instance.MarkAsRead(db.DbConnection, notificationInstance)
	request.Write(notificationInstance)
}

func (self *NotificationHandler) Post(request *gottp.Request) {
	notificationInstance := new(notification_instance.NotificationInstance)
	request.ConvertArguments(notificationInstance)
	notificationInstance.NotificationType = request.GetArgument("generic_id").(string)
	notificationInstance.IsRead = false

	notificationInstance.PrepareSave()

	if !utils.ValidateAndRaiseError(request, notificationInstance) {
		return
	}

	if !notificationInstance.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Context is required while inserting."})
		return
	}

	notificationSetting := notification.FindAppropriateNotification(db.DbConnection, notificationInstance.User, notificationInstance.ApplicationID, notificationInstance.Organization, notificationInstance.NotificationType)

	if notificationSetting == nil {
		log.Println("Unable to find suitable notification setting.")
		return
	} else {
		notifications.SendNotification(db.DbConnection, notificationInstance, notificationSetting)
	}

	notification_instance.InsertIntoDatabase(db.DbConnection, notificationInstance)

}
