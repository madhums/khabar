package models

import (
	"github.com/parthdesai/sc-notifications/db"
)

const (
	UserLocaleCollection = "user_locales"
)

type UserLocale struct {
	db.BaseModel `bson:",inline"`
	UserId       string `json:"user_id" bson:"user_id" required:"true"`
	LanguageId   string `json:"language_id" bson:"language_id" required:"true"`
	RegionId     string `json:"region_id" bson:"region_id" required:"true"`
}

func (self *UserLocale) GetFromDatabase(dbConn *db.MConn) bool {
	return dbConn.Get(UserLocaleCollection, db.M{"user_id": self.UserId}).Next(self)
}

func (self *UserLocale) InsertIntoDatabase(dbConn *db.MConn) string {
	return dbConn.Insert(UserLocaleCollection, self)
}

func (self *UserLocale) IsValid() bool {
	if len(self.LanguageId) == 0 || len(self.RegionId) == 0 || len(self.UserId) == 0 {
		return false
	}
	return true
}

func (self *UserLocale) Update(dbConn *db.MConn) error {
	return dbConn.Update(UserLocaleCollection, db.M{"_id": self.Id},
		db.M{
			"$set": db.M{
				"region_id":   self.RegionId,
				"language_id": self.LanguageId,
			},
		})

}
