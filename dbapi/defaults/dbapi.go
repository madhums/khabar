/**
 * What does this migration do?
 *
 * - Removes app_name from topics collection if it exists
 * - Modifies `channels` property in `topics` collection from array of strings
 *   to array of objects
 * - Removes previously added `value` property to `topics` collection
 */

package defaults

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2/bson"
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

func IsDefaultExists(org string, topic, channel string, enabled bool) bool {

	return db.Conn.Count(db.DefaultsCollection, utils.M{
		"org":      org,
		"ident":    topic,
		"channels": channel,
		"enabled":  enabled,
	}) != 0
}

func Get(org, topic string, enabled bool) (error, *db.Defaults) {
	defaultPref := new(db.Defaults)

	err := db.Conn.GetOne(db.DefaultsCollection, utils.M{
		"org":     org,
		"ident":   topic,
		"enabled": enabled,
	}, defaultPref)

	if err != nil {
		return err, nil
	}

	return nil, defaultPref
}

func Update(id bson.ObjectId, doc utils.M) error {
	return db.Conn.Update(db.DefaultsCollection, utils.M{"_id": id}, doc)
}

func GetAllEnabled(org string) []db.Defaults {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Defaults{}

	db.Conn.Get(session, db.DefaultsCollection, utils.M{
		"org":     org,
		"enabled": true,
	}).All(&result)

	return result
}

func GetAllDisabled(org string) []db.Defaults {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Defaults{}

	db.Conn.Get(session, db.DefaultsCollection, utils.M{
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

	_, err := db.Conn.FindAndUpdate(db.DefaultsCollection, query, spec, &result)
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

	_, err := db.Conn.FindAndUpdate(db.DefaultsCollection, query, spec, &result)
	return err
}
