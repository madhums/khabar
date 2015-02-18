package notification_instance

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v1"
)

func Update(dbConn *db.MConn, id bson.ObjectId, doc *db.M) error {
	return dbConn.Update(NotificationInstanceCollection, db.M{
		"_id": id,
	}, db.M{
		"$set": *doc,
	})
}

func GetAll(dbConn *db.MConn, paginator *gottp.Paginator, user string, appName string, organization string) *[]NotificationInstance {
	var query db.M = nil
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []NotificationInstance
	query["user"] = user

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(organization) > 0 {
		query["org"] = organization
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

func Insert(dbConn *db.MConn, notificationInstance *NotificationInstance) string {
	return dbConn.Insert(NotificationInstanceCollection, notificationInstance)
}
