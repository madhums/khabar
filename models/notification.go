package models

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/utils"
)

type Notification struct {
	db.BaseModel   `bson:",inline"`
	UserID         string   `json:"user_id" bson:"user_id"`
	OrganizationID string   `json:"org_id" bson:"org_id"`
	ApplicationID  string   `json:"app_id" bson:"app_id"`
	Channels       []string `json:"channels" bson:"channels" required:"true"`
	Type           string   `json:"type" bson:"type" required:"true"`
}

func (self *Notification) IsValid(op_type int) bool {
	if (len(self.UserID) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Type) == 0 {
		return false
	}

	if op_type == INSERT_OPERATION {

		if len(self.Channels) == 0 {
			return false
		}
	}

	return true
}

func (self *Notification) UpdateChannels(dbConn *db.MConn) error {
	return db.DbConnection.Update("notifications", db.M{"_id": self.Id},
		db.M{
			"$set": db.M{
				"channels": self.Channels,
			},
		})
}

func (self *Notification) DeleteFromDatabase(dbConn *db.MConn) error {
	return dbConn.Delete("notifications", db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "type": self.Type})
}

func (self *Notification) GetFromDatabase(dbConn *db.MConn) bool {
	return dbConn.Get("notifications", db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "type": self.Type}).Next(self)
}

func (self *Notification) AddChannelToNotification(channel string) {
	newArray := make([]string, len(self.Channels)+1)
	copy(newArray, self.Channels)
	newArray[len(newArray)-1] = channel
	self.Channels = newArray

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
