package notification

import (
	"github.com/parthdesai/sc-notifications/db"
)

func Update(dbConn *db.MConn, user string, appName string, organization string, notificationType string, doc *db.M) error {

	return dbConn.Update(NotificationCollection,
		db.M{"app_name": appName,
			"org":  organization,
			"user": user,
			"type": notificationType,
		},
		db.M{
			"$set": *doc,
		})
}

func Insert(dbConn *db.MConn, notification *Notification) string {
	return dbConn.Insert(NotificationCollection, notification)
}

func Delete(dbConn *db.MConn, doc *db.M) error {
	return dbConn.Delete(NotificationCollection, *doc)
}

func Get(dbConn *db.MConn, user string, appName string, organization string, notificationType string) *Notification {
	notification := new(Notification)
	if dbConn.GetOne(NotificationCollection, db.M{"app_name": appName,
		"org": organization, "user": user, "type": notificationType}, notification) != nil {
		return nil
	}
	return notification
}

func FindAppropriateNotificationForUser(dbConn *db.MConn, user string, appName string, organization string, notificationType string) *Notification {
	var err error
	notification := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     user,
		"app_name": appName,
		"org":      organization,
		"type":     notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":     user,
		"app_name": appName,
		"type":     notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user": user,
		"org":  organization,
		"type": notificationType,
	}, notification)

	if err == nil {
		return notification
	}
	return nil
}

func FindAppropriateOrganizationNotification(dbConn *db.MConn, appName string, organization string, notificationType string) *Notification {
	var err error
	notification := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"app_name": appName,
		"org":      organization,
		"type":     notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"org":  organization,
		"type": notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	return nil

}

func FindGlobalNotification(dbConn *db.MConn, notificationType string) *Notification {
	var err error
	notification := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"type": notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	return nil

}

func FindAppropriateNotification(dbConn *db.MConn, user string, appName string, organization string, notificationType string) *Notification {
	var notification *Notification

	notification = FindAppropriateNotificationForUser(dbConn, user, appName, organization, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindAppropriateOrganizationNotification(dbConn, appName, organization, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindGlobalNotification(dbConn, notificationType)

	if notification != nil {
		return notification
	}

	return nil

}
