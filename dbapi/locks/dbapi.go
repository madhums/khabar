package locks

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
)

func Insert(preference *db.Locks) string {
	return db.Conn.Insert(db.LocksCollection, preference)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.LocksCollection, *doc)
}

func GetAll(user, org string) []db.Locks {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Locks{}

	db.Conn.Get(session, db.LocksCollection, utils.M{
		"user": user,
		"org":  org,
	}).All(result)

	return result
}

func GetAllEnabled(user, org string) []db.Locks {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Locks{}

	db.Conn.Get(session, db.LocksCollection, utils.M{
		"user":    user,
		"org":     org,
		"enabled": true,
	}).All(result)

	return result
}

func GetAllDisabled(user, org string) []db.Locks {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Locks{}

	db.Conn.Get(session, db.LocksCollection, utils.M{
		"user":    user,
		"org":     org,
		"enabled": false,
	}).All(result)

	return result
}
