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
	LanguageID   string `json:"language_id" bson:"language_id" required:"true"`
	RegionID     string `json:"region_id" bson:"region_id" required:"true"`
}

func (self *UserLocale) IsValid() bool {
	if len(self.LanguageID) == 0 || len(self.RegionID) == 0 || len(self.UserID) == 0 {
		return false
	}
	return true
}
