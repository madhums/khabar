package notification_instance

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
)

func MarkAsRead(dbConn *db.MConn, notificationInstance *NotificationInstance) {
	dbConn.FindAndUpdate(NotificationInstanceCollection,
		db.M{
			"_id": notificationInstance.Id,
		},
		db.M{
			"$set": db.M{
				"is_read": true,
			},
		}, notificationInstance)
}

func GetAllFromDatabase(dbConn *db.MConn, paginator *gottp.Paginator, user string, applicationID string, organizationID string) *[]NotificationInstance {
	var query db.M = nil
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []NotificationInstance
	query["user"] = user

	if len(applicationID) > 0 {
		query["app_id"] = applicationID
	}

	if len(organizationID) > 0 {
		query["org_id"] = organizationID
	}

	var limit int
	var skip int
	var limitExists bool
	var skipExists bool

	limit, limitExists = query["limit"].(int)
	skip, skipExists = query["skip"].(int)

	delete(query, "limit")
	delete(query, "skip")

	if !limitExists {
		limit = 30
	}
	if !skipExists {
		skip = 0
	}

	delete(query, "limit")
	delete(query, "skip")

	dbConn.GetCursor(NotificationInstanceCollection, query).Skip(skip).Limit(limit).All(&result)

	return &result
}

func InsertIntoDatabase(dbConn *db.MConn, notificationInstance *NotificationInstance) string {
	return dbConn.Insert(NotificationInstanceCollection, notificationInstance)
}
