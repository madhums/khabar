package handlers

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi/notification"
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"github.com/changer/sc-notifications/notifications"
	"github.com/changer/sc-notifications/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v1"
	"log"
	"net/http"
)

type Notification struct {
	gottp.BaseHandler
}

func (self *Notification) Get(request *gottp.Request) {

	notificationInstance := new(notification_instance.NotificationInstance)
	request.ConvertArguments(notificationInstance)

	notificationInstance.User = request.GetArgument("generic_id").(string)

	paginator := request.GetPaginator()

	request.Write(notification_instance.GetAll(db.Conn, paginator, notificationInstance.User, notificationInstance.AppName, notificationInstance.Organization))
}

func (self *Notification) Put(request *gottp.Request) {
	notificationInstance := new(notification_instance.NotificationInstance)
	objectIdString := request.GetArgument("generic_id").(string)
	if !bson.IsObjectIdHex(objectIdString) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "_id is not a valid Hex object."})
		return
	}
	notificationInstance.Id = bson.ObjectIdHex(objectIdString)
	notification_instance.Update(db.Conn, notificationInstance.Id, &db.M{
		"is_read": true,
	})
	notificationInstance.IsRead = true
	request.Write(notificationInstance)
}

func (self *Notification) Post(request *gottp.Request) {
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

	notificationSetting := notification.FindAppropriateNotification(db.Conn, notificationInstance.User, notificationInstance.AppName, notificationInstance.Organization, notificationInstance.NotificationType)

	if notificationSetting == nil {
		log.Println("Unable to find suitable notification setting.")
		return
	} else {
		notifications.SendNotification(db.Conn, notificationInstance, notificationSetting)
	}

	notification_instance.Insert(db.Conn, notificationInstance)

}
