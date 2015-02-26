package handlers

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi/notification"
	"github.com/changer/sc-notifications/dbapi/notification_instance"
	"github.com/changer/sc-notifications/dbapi/sent_notification"
	"github.com/changer/sc-notifications/notifications"
	"github.com/changer/sc-notifications/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v2"
	"log"
	"net/http"
)

type Notification struct {
	gottp.BaseHandler
}

func (self *Notification) Get(request *gottp.Request) {

	ntfInst := new(sent_notification.NotificationInstance)
	request.ConvertArguments(ntfInst)

	ntfInst.User = request.GetArgument("generic_id").(string)

	paginator := request.GetPaginator()

	request.Write(sent_notification.GetAll(db.Conn, paginator, ntfInst.User, ntfInst.AppName, ntfInst.Organization))
}

func (self *Notification) Put(request *gottp.Request) {
	ntfInst := new(sent_notification.NotificationInstance)
	objectIdString := request.GetArgument("generic_id").(string)
	if !bson.IsObjectIdHex(objectIdString) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "_id is not a valid Hex object."})
		return
	}
	ntfInst.Id = bson.ObjectIdHex(objectIdString)
	sent_notification.Update(db.Conn, ntfInst.Id, &db.M{
		"is_read": true,
	})
	ntfInst.IsRead = true
	request.Write(ntfInst)
}

func (self *Notification) Post(request *gottp.Request) {
	ntfInst := new(notification_instance.NotificationInstance)
	request.ConvertArguments(ntfInst)
	ntfInst.NotificationType = request.GetArgument("generic_id").(string)
	ntfInst.IsRead = false

	ntfInst.PrepareSave()

	if !utils.ValidateAndRaiseError(request, ntfInst) {
		return
	}

	if !ntfInst.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Context is required while inserting."})
		return
	}

	ntfSetting := notification.FindAppropriateNotification(db.Conn, ntfInst.User, ntfInst.AppName, ntfInst.Organization, ntfInst.NotificationType)

	if ntfSetting == nil {
		log.Println("Unable to find suitable notification setting.")
		return
	} else {
		notifications.SendNotification(db.Conn, ntfInst, ntfSetting)
	}

}
