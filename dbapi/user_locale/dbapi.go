package user_locale

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

func Get(dbConn *db.MConn, user string) *UserLocale {
	userLocale := new(UserLocale)
	if dbConn.GetOne(UserLocaleCollection, utils.M{"user": user}, userLocale) != nil {
		return nil
	}
	return userLocale
}

func Insert(dbConn *db.MConn, userLocale *UserLocale) string {
	return dbConn.Insert(UserLocaleCollection, userLocale)
}

func Update(dbConn *db.MConn, user string, doc *utils.M) error {
	return dbConn.Update(UserLocaleCollection, utils.M{"user": user},
		utils.M{
			"$set": *doc,
		})
}
