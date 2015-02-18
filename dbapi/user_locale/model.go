package user_locale

import (
	"github.com/changer/sc-notifications/db"
)

const (
	UserLocaleCollection = "user_locales"
)

type UserLocale struct {
	db.BaseModel `bson:",inline"`
	User         string `json:"user" bson:"user" required:"true"`
	Locale       string `json:"locale" bson:"locale" required:"true"`
	TimeZone     string `json:"timezone" bson:"timezone" required:"true"`
}

func (self *UserLocale) IsValid() bool {
	if len(self.Locale) == 0 || len(self.User) == 0 || len(self.TimeZone) == 0 {
		return false
	}
	return true
}
