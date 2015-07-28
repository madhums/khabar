package locks

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2/bson"
)

func Insert(preference *db.Locks) string {
	return db.Conn.Insert(db.LocksCollection, preference)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.LocksCollection, *doc)
}

func GetAll(org string) []db.Topic {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Topic{}

	db.Conn.Get(session, db.TopicCollection, utils.M{
		"$or": utils.M{
			"org":             org,
			"user":            "",
			"channels.locked": true,
		},
	}).All(&result)

	return result
}

func Get(org, topic string, enabled bool) (error, *db.Locks) {
	lock := new(db.Locks)

	err := db.Conn.GetOne(db.LocksCollection, utils.M{
		"org":     org,
		"ident":   topic,
		"enabled": enabled,
	}, lock)

	if err != nil {
		return err, nil
	}

	return nil, lock
}

func Update(id bson.ObjectId, doc utils.M) error {
	return db.Conn.Update(db.LocksCollection, utils.M{"_id": id}, doc)
}

func IsLocked(org string, topic, channel string, enabled bool) bool {
	return db.Conn.Count(db.LocksCollection, utils.M{
		"org":      org,
		"ident":    topic,
		"channels": channel,
		"enabled":  enabled,
	}) != 0
}

func GetAllEnabled(org string) []db.Locks {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Locks{}

	db.Conn.Get(session, db.LocksCollection, utils.M{
		"org":     org,
		"enabled": true,
	}).All(&result)

	return result
}

func GetAllDisabled(org string) []db.Locks {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Locks{}

	db.Conn.Get(session, db.LocksCollection, utils.M{
		"org":     org,
		"enabled": false,
	}).All(&result)

	return result
}

func AddChannel(topic, channel, organization string, enabled bool) error {

	query := utils.M{
		"org":     organization,
		"ident":   topic,
		"enabled": enabled,
	}

	spec := utils.M{"$addToSet": utils.M{"channels": channel}}

	result := utils.M{}

	_, err := db.Conn.FindAndUpdate(db.LocksCollection, query, spec, &result)
	return err
}

func RemoveChannel(topic, channel, organization string, enabled bool) error {

	query := utils.M{
		"org":     organization,
		"ident":   topic,
		"enabled": enabled,
	}

	spec := utils.M{"$pull": utils.M{"channels": channel}}

	result := utils.M{}

	_, err := db.Conn.FindAndUpdate(db.LocksCollection, query, spec, &result)
	return err
}
