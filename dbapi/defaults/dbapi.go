package defaults

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
)

func Insert(preference *db.Defaults) string {
	return db.Conn.Insert(db.DefaultsCollection, preference)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.DefaultsCollection, *doc)
}

func GetAll(org string) []db.Defaults {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Defaults{}

	db.Conn.Get(session, db.DefaultsCollection, utils.M{
		"org": org,
	}).All(&result)

	return result
}

func GetAllEnabled(org string) []db.Defaults {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Defaults{}

	db.Conn.Get(session, db.DefaultsCollection, utils.M{
		"org":     org,
		"enabled": "true",
	}).All(&result)

	return result
}

func GetAllDisabled(org string) []db.Defaults {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Defaults{}

	db.Conn.Get(session, db.DefaultsCollection, utils.M{
		"org":     org,
		"enabled": "false",
	}).All(&result)

	return result
}
