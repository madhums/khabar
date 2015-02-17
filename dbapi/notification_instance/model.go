package notification_instance

import (
	"github.com/parthdesai/sc-notifications/db"
)

const (
	NotificationInstanceCollection = "notification_instances"
)

type NotificationInstance struct {
	db.BaseModel     `bson:",inline"`
	Organization     string                 `json:"org" bson:"org" required:"true"`
	AppName          string                 `json:"app_name" bson:"app_name" required:"true"`
	NotificationType string                 `json:"notification_type" bson:"notification_type" required:"true"`
	IsPending        bool                   `json:"is_pending" bson:"is_pending" required:"true"`
	User             string                 `json:"user" bson:"user" required:"true"`
	Context          map[string]interface{} `json:"context" bson:"context" required:"true"`
	IsRead           bool                   `json:"is_read" bson:"is_read"`
}

func (self *NotificationInstance) IsValid() bool {
	if len(self.Context) == 0 {
		return false
	}
	return true
}
