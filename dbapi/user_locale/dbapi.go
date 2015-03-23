package user_locale

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

func Get(dbConn *db.MConn, user string) (userLocale *UserLocale, err error) {
	userLocale = new(UserLocale)
	err = dbConn.GetOne(UserLocaleCollection, utils.M{"user": user}, userLocale)
	return
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
