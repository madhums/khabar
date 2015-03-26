package topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

type Topic struct {
	db.BaseModel `bson:",inline"`
	User         string   `json:"user" bson:"user"`
	Organization string   `json:"org" bson:"org"`
	AppName      string   `json:"app_name" bson:"app_name"`
	Channels     []string `json:"channels" bson:"channels" required:"true"`
	Ident        string   `json:"ident" bson:"ident" required:"true"`
}

func (self *Topic) IsValid(op_type int) bool {
	if (len(self.User) == 0) && (len(self.Organization) == 0) &&
		(len(self.AppName) == 0) {
		return false
	}

	if len(self.Ident) == 0 {
		return false
	}

	if op_type == db.INSERT_OPERATION {

		if len(self.Channels) == 0 {
			return false
		}
	}

	return true
}

func (self *Topic) AddChannel(channel string) {
	self.Channels = append(self.Channels, channel)
	utils.RemoveDuplicates(&(self.Channels))
}

func (self *Topic) RemoveChannel(channel string) {
	j := 0
	for i, x := range self.Channels {
		if x != channel {
			self.Channels[j] = self.Channels[i]
			j++
		}
	}
	self.Channels = self.Channels[:j]
}
