package sent

import (
	"github.com/changer/khabar/db"
)

const (
	SentCollection = "sent_notifications"
)

type SentItem struct {
	db.BaseModel   `bson:",inline"`
	Organization   string `json:"org" bson:"org" required:"true"`
	AppName        string `json:"app_name" bson:"app_name" required:"true"`
	Topic          string `json:"topic" bson:"topic" required:"true"`
	User           string `json:"user" bson:"user" required:"true"`
	DestinationUri string `json:"destination_uri" bson:"destination_uri" required:"true"`
	Text           string `json:"text" bson:"text" required:"true"`
	IsRead         bool   `json:"is_read" bson:"is_read"`
}

func (self *SentItem) IsValid() bool {
	if len(self.Text) == 0 {
		return false
	}
	return true
}
