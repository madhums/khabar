package models

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

func (self *UserLocale) GetFromDatabase(dbConn *db.MConn) bool {
	return dbConn.Get(UserLocaleCollection, db.M{"user_id": self.UserID}).Next(self)
}

func (self *UserLocale) InsertIntoDatabase(dbConn *db.MConn) string {
	return dbConn.Insert(UserLocaleCollection, self)
}

func (self *UserLocale) IsValid() bool {
	if len(self.LanguageID) == 0 || len(self.RegionID) == 0 || len(self.UserID) == 0 {
		return false
	}
	return true
}

func (self *UserLocale) Update(dbConn *db.MConn) error {
	return dbConn.Update(UserLocaleCollection, db.M{"_id": self.Id},
		db.M{
			"$set": db.M{
				"region_id":   self.RegionID,
				"language_id": self.LanguageID,
			},
		})

}
