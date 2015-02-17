package user_locale

import (
	"github.com/parthdesai/sc-notifications/db"
)

func GetFromDatabase(dbConn *db.MConn, userID string) *UserLocale {
	userLocale := new(UserLocale)
	if dbConn.GetOne(UserLocaleCollection, db.M{"user_id": userID}, userLocale) != nil {
		return nil
	}
	return userLocale
}

func InsertIntoDatabase(dbConn *db.MConn, userLocale *UserLocale) string {
	return dbConn.Insert(UserLocaleCollection, userLocale)
}

func Update(dbConn *db.MConn, userLocale *UserLocale) error {
	return dbConn.Update(UserLocaleCollection, db.M{"_id": userLocale.Id},
		db.M{
			"$set": db.M{
				"locale":   userLocale.Locale,
				"timezone": userLocale.TimeZone,
			},
		})

}
