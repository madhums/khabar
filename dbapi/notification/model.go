package notification

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi"
	"github.com/parthdesai/sc-notifications/utils"
	"log"
)

const (
	NotificationCollection = "notifications"
)

type Notification struct {
	db.BaseModel   `bson:",inline"`
	User           string   `json:"user" bson:"user"`
	OrganizationID string   `json:"org_id" bson:"org_id"`
	ApplicationID  string   `json:"app_id" bson:"app_id"`
	Channels       []string `json:"channels" bson:"channels" required:"true"`
	Type           string   `json:"type" bson:"type" required:"true"`
}

func (self *Notification) IsValid(op_type int) bool {
	if (len(self.User) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Type) == 0 {
		return false
	}

	if op_type == dbapi.INSERT_OPERATION {

		if len(self.Channels) == 0 {
			return false
		}
	}

	return true
}

func (self *Notification) AddChannelToNotification(channel string) {
	log.Println(self.Channels)
	self.Channels = append(self.Channels, channel)
	log.Println(self.Channels)
	utils.RemoveDuplicates(&(self.Channels))
}

func (self *Notification) RemoveChannelFromNotification(channel string) {
	j := 0
	for i, x := range self.Channels {
		if x != channel {
			self.Channels[j] = self.Channels[i]
			j++
		}
	}
	self.Channels = self.Channels[:j]
}
