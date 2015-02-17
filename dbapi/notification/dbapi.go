package notification

import (
	"github.com/parthdesai/sc-notifications/db"
)

func UpdateChannels(dbConn *db.MConn, notification *Notification) error {
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
		"org_id": notification.OrganizationID, "user_id": notification.UserID, "type": notification.Type})
}

func GetFromDatabase(dbConn *db.MConn, userID string, applicationID string, organizationID string, notificationType string) *Notification {
	notification := new(Notification)
	if !dbConn.Get(NotificationCollection, db.M{"app_id": applicationID,
		"org_id": organizationID, "user_id": userID, "type": notificationType}).Next(notification) {
		return nil
	}
	return notification
}

func FindAppropriateNotificationForUser(dbConn *db.MConn, userID string, applicationID string, organizationID string, notificationType string) *Notification {
	var hasData bool
	notification := new(Notification)
	hasData = dbConn.Get(NotificationCollection, db.M{
		"user_id": userID,
		"app_id":  applicationID,
		"org_id":  organizationID,
		"type":    notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}

	hasData = dbConn.Get(NotificationCollection, db.M{
		"user_id": userID,
		"app_id":  applicationID,
		"type":    notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}

	hasData = dbConn.Get(NotificationCollection, db.M{
		"user_id": userID,
		"org_id":  organizationID,
		"type":    notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}
	return nil
}

func FindAppropriateOrganizationNotification(dbConn *db.MConn, applicationID string, organizationID string, notificationType string) *Notification {
	var hasData bool
	notification := new(Notification)
	hasData = dbConn.Get(NotificationCollection, db.M{
		"app_id": applicationID,
		"org_id": organizationID,
		"type":   notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}

	hasData = dbConn.Get(NotificationCollection, db.M{
		"org_id": organizationID,
		"type":   notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}

	return nil

}

func FindGlobalNotification(dbConn *db.MConn, notificationType string) *Notification {
	var hasData bool
	notification := new(Notification)
	hasData = dbConn.Get(NotificationCollection, db.M{
		"type": notificationType,
	}).Next(notification)

	if hasData {
		return notification
	}

	return nil

}

func FindAppropriateNotification(dbConn *db.MConn, userID string, applicationID string, organizationID string, notificationType string) *Notification {
	var notification *Notification

	notification = FindAppropriateNotificationForUser(dbConn, userID, applicationID, organizationID, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindAppropriateOrganizationNotification(dbConn, applicationID, organizationID, notificationType)

	if notification != nil {
		return notification
	}

	notification = FindGlobalNotification(dbConn, notificationType)

	if notification != nil {
		return notification
	}

	return nil

}
