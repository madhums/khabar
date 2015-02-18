package user_locale

import (
	"github.com/changer/sc-notifications/db"
)

func Get(dbConn *db.MConn, user string) *UserLocale {
	userLocale := new(UserLocale)
	if dbConn.GetOne(UserLocaleCollection, db.M{"user": user}, userLocale) != nil {
		return nil
	}
	return userLocale
}

func Insert(dbConn *db.MConn, userLocale *UserLocale) string {
	return dbConn.Insert(UserLocaleCollection, userLocale)
}

func Update(dbConn *db.MConn, user string, doc *db.M) error {
	return dbConn.Update(UserLocaleCollection, db.M{"user": user},
		db.M{
			"$set": *doc,
		})
}
