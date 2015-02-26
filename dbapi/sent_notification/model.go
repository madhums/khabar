package sent_notification

import (
	"github.com/changer/sc-notifications/db"
)

const (
	NotificationInstanceCollection = "sent_notifications"
)

type NotificationInstance struct {
	db.BaseModel     `bson:",inline"`
	Organization     string `json:"org" bson:"org" required:"true"`
	AppName          string `json:"app_name" bson:"app_name" required:"true"`
	NotificationType string `json:"notification_type" bson:"notification_type" required:"true"`
	User             string `json:"user" bson:"user" required:"true"`
	DestinationUri   string `json:"destination_uri" bson:"destination_uri" required:"true"`
	NotificationText string `json:"notification_text" bson:"notification_text" required:"true"`
	IsRead           bool   `json:"is_read" bson:"is_read"`
}

func (self *NotificationInstance) IsValid() bool {
	if len(self.NotificationText) == 0 {
		return false
	}
	return true
}
