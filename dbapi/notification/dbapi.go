package notification

import (
	"github.com/changer/sc-notifications/db"
)

func Update(dbConn *db.MConn, user string, appName string, org string, ntfType string, doc *db.M) error {

	return dbConn.Update(NotificationCollection,
		db.M{"app_name": appName,
			"org":  org,
			"user": user,
			"type": ntfType,
		},
		db.M{
			"$set": *doc,
		})
}

func Insert(dbConn *db.MConn, ntf *Notification) string {
	return dbConn.Insert(NotificationCollection, ntf)
}

func Delete(dbConn *db.MConn, doc *db.M) error {
	return dbConn.Delete(NotificationCollection, *doc)
}

func Get(dbConn *db.MConn, user string, appName string, org string, ntfType string) *Notification {
	ntf := new(Notification)
	if dbConn.GetOne(NotificationCollection, db.M{"app_name": appName,
		"org": org, "user": user, "type": ntfType}, ntf) != nil {
		return nil
	}
	return ntf
}

func FindAppropriateNotificationForUser(dbConn *db.MConn, user string, appName string, org string, ntfType string) *Notification {
	var err error
	ntf := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     user,
		"app_name": appName,
		"org":      org,
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     user,
		"app_name": appName,
		"org":      "",
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     user,
		"app_name": "",
		"org":      org,
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}
	return nil
}

func FindAppropriateOrganizationNotification(dbConn *db.MConn, appName string, org string, ntfType string) *Notification {
	var err error
	ntf := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     "",
		"app_name": appName,
		"org":      org,
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     "",
		"app_name": "",
		"org":      org,
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}

	return nil

}

func FindGlobalNotification(dbConn *db.MConn, ntfType string) *Notification {
	var err error
	ntf := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     "",
		"app_name": "",
		"org":      "",
		"type":     ntfType,
	}, ntf)

	if err == nil {
		return ntf
	}

	return nil

}

func FindAppropriateNotification(dbConn *db.MConn, user string, appName string, org string, ntfType string) *Notification {
	var ntf *Notification

	ntf = FindAppropriateNotificationForUser(dbConn, user, appName, org, ntfType)

	if ntf != nil {
		return ntf
	}

	ntf = FindAppropriateOrganizationNotification(dbConn, appName, org, ntfType)

	if ntf != nil {
		return ntf
	}

	ntf = FindGlobalNotification(dbConn, ntfType)

	if ntf != nil {
		return ntf
	}

	return nil

}
