package user_locale

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

func Get(dbConn *db.MConn, user string) (userLocale *db.UserLocale, err error) {
	userLocale = new(db.UserLocale)
	err = dbConn.GetOne(db.UserLocaleCollection, utils.M{"user": user}, userLocale)
	return
}

func Insert(dbConn *db.MConn, userLocale *db.UserLocale) string {
	return dbConn.Insert(db.UserLocaleCollection, userLocale)
}

func Update(dbConn *db.MConn, user string, doc *utils.M) error {
	return dbConn.Update(db.UserLocaleCollection, utils.M{"user": user},
		utils.M{
			"$set": *doc,
		})
}
