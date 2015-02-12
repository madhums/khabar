package models

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
)

const (
	NotificationInstanceCollection = "notification_instances"
)

type NotificationInstance struct {
	db.BaseModel     `bson:",inline"`
	OrganizationID   string                 `json:"org_id" bson:"org_id" required:"true"`
	ApplicationID    string                 `json:"app_name" bson:"app_name" required:"true"`
	NotificationType string                 `json:"notification_type" bson:"notification_type" required:"true"`
	IsPending        bool                   `json:"is_pending" bson:"is_pending" required:"true"`
	UserID           string                 `json:"user_id" bson:"user_id" required:"true"`
	Context          map[string]interface{} `json:"context" bson:"context" required:"true"`
	IsRead           bool                   `json:"is_read" bson:"is_read"`
}

func (self *NotificationInstance) IsValid() bool {
	if len(self.Context) == 0 {
		return false
	}
	return true
}

func (self *NotificationInstance) MarkAsRead(dbConn *db.MConn) {
	dbConn.FindAndUpdate(NotificationInstanceCollection,
		db.M{
			"_id": self.Id,
		},
		db.M{
			"$set": db.M{
				"is_read": true,
			},
		},
		self)
}

func (self *NotificationInstance) GetAllFromDatabase(dbConn *db.MConn, paginator *gottp.Paginator) *[]NotificationInstance {
	var query db.M = nil
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []NotificationInstance
	query["user_id"] = self.UserID

	if len(self.ApplicationID) > 0 {
		query["app_id"] = self.ApplicationID
	}

	if len(self.OrganizationID) > 0 {
		query["org_id"] = self.OrganizationID
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

func (self *NotificationInstance) InsertIntoDatabase(dbConn *db.MConn) string {
	return dbConn.Insert(NotificationInstanceCollection, self)
}
