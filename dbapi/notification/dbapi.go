package notification

import (
	"github.com/parthdesai/sc-notifications/db"
)

func UpdateNotification(dbConn *db.MConn, notification *Notification) error {
	return dbConn.FindAndUpdate(NotificationCollection, db.M{"_id": notification.Id},
		db.M{
			"$set": db.M{
				"channels": notification.Channels,
			},
		}, notification)
}

func InsertIntoDatabase(dbConn *db.MConn, notification *Notification) string {
	return dbConn.Insert(NotificationCollection, notification)
}

func DeleteFromDatabase(dbConn *db.MConn, notification *Notification) error {
	return dbConn.Delete(NotificationCollection, db.M{"app_id": notification.ApplicationID,
		"org": notification.Organization, "user": notification.User, "type": notification.Type})
}

func GetFromDatabase(dbConn *db.MConn, user string, applicationID string, organization string, notificationType string) *Notification {
	notification := new(Notification)
	if dbConn.GetOne(NotificationCollection, db.M{"app_id": applicationID,
		"org": organization, "user": user, "type": notificationType}, notification) != nil {
		return nil
	}
	return notification
}

func FindAppropriateNotificationForUser(dbConn *db.MConn, user string, applicationID string, organization string, notificationType string) *Notification {
	var err error
	notification := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":   user,
		"app_id": applicationID,
		"org":    organization,
		"type":   notificationType,
	}, notification)

	if err == nil {
		return notification
	}

	err = dbConn.GetOne(NotificationCollection, db.M{
		"user":   user,
		"app_id": applicationID,
		"type":   notificationType,
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

func FindAppropriateOrganizationNotification(dbConn *db.MConn, applicationID string, organization string, notificationType string) *Notification {
	var err error
	notification := new(Notification)
	err = dbConn.GetOne(NotificationCollection, db.M{
		"app_id": applicationID,
		"org":    organization,
		"type":   notificationType,
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

func FindAppropriateNotification(dbConn *db.MConn, user string, applicationID string, organization string, notificationType string) *Notification {
	var notification *Notification

	notification = FindAppropriateNotificationForUser(dbConn, user, applicationID, organization, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindAppropriateOrganizationNotification(dbConn, applicationID, organization, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindGlobalNotification(dbConn, notificationType)

	if notification != nil {
		return notification
	}

	return nil

}
