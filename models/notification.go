package models

import "github.com/parthdesai/sc-notifications/db"

type Notification struct {
	db.BaseModel   `bson:",inline"`
	UserID         string   `json:"user_id" bson:"user_id"`
	OrganizationID string   `json:"org_id" bson:"org_id"`
	ApplicationID  string   `json:"app_id" bson:"app_id"`
	Channels       []string `json:"channels" bson:"channels" required:"true"`
	Type           string   `json:"type" bson:"type" required:"true"`
}

func (self *Notification) IsValid() bool {
	if (len(self.UserID) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Type) == 0 {
		return false
	}

	if len(self.Channels) == 0 {
		return false
	}

	return true
}
