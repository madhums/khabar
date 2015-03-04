package pending

import (
	"github.com/changer/khabar/db"
)

type PendingItem struct {
	db.BaseModel   `bson:",inline"`
	Organization   string                 `json:"org" bson:"org" required:"true"`
	AppName        string                 `json:"app_name" bson:"app_name" required:"true"`
	Topic          string                 `json:"topic" bson:"topic" required:"true"`
	IsPending      bool                   `json:"is_pending" bson:"is_pending" required:"true"`
	User           string                 `json:"user" bson:"user" required:"true"`
	DestinationUri string                 `json:"destination_uri" bson:"destination_uri" required:"true"`
	Context        map[string]interface{} `json:"context" bson:"context" required:"true"`
	IsRead         bool                   `json:"is_read" bson:"is_read"`
}

func (self *PendingItem) IsValid() bool {
	if len(self.Context) == 0 {
		return false
	}
	return true
}
