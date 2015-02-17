package user_locale

import (
	"github.com/parthdesai/sc-notifications/db"
)

const (
	UserLocaleCollection = "user_locales"
)

type UserLocale struct {
	db.BaseModel `bson:",inline"`
	UserID       string `json:"user_id" bson:"user_id" required:"true"`
	Locale       string `json:"locale" bson:"locale" required:"true"`
	TimeZone     string `json:"timezone" bson:"timezone" required:"true"`
}

func (self *UserLocale) IsValid() bool {
	if len(self.Locale) == 0 || len(self.UserID) == 0 || len(self.TimeZone) == 0 {
		return false
	}
	return true
}
