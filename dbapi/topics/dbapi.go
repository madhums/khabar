package topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
)

const BLANK = ""

func Update(user string, appName string,
	org string, topicName string, doc *utils.M) error {

	return db.Conn.Update(db.TopicCollection,
		utils.M{"app_name": appName,
			"org":   org,
			"user":  user,
			"ident": topicName,
		},
		utils.M{
			"$set": *doc,
		})
}

func Insert(topic *Topic) string {
	return db.Conn.Insert(db.TopicCollection, topic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.TopicCollection, *doc)
}

func ChannelAllowed(user, appname, org, topicName, channel string) bool {
	return db.Conn.Count(db.TopicCollection, utils.M{
		"$or": []utils.M{
			utils.M{"org": org},
			utils.M{"app_name": appname},
			utils.M{"user": user},
		},
		"ident":    topicName,
		"channels": channel,
	}) == 0
}

func Get(user, appName, org,
	topicName string) (topic *Topic, err error) {

	topic = new(Topic)

	err = db.Conn.GetOne(
		db.TopicCollection,
		utils.M{
			"app_name": appName,
			"org":      org,
			"user":     user,
			"ident":    topicName,
		},
		topic,
	)

	if err != nil {
		return nil, err
	}

	return
}

func GetAll(user, appName, org string) (*mgo.Iter, error) {
	var query utils.M = make(utils.M)

	if len(user) > 0 {
		query["user"] = user
	}

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	session := db.Conn.Session.Copy()
	defer session.Close()

	iter := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return iter, nil
}

func findPerUser(user, appName, org,
	topicName string) (topic *Topic, err error) {

	topic, err = Get(user, appName, org, topicName)
	if err != nil {
		topic, err = Get(user, appName, BLANK, topicName)
		if err != nil {
			topic, err = Get(user, BLANK, org, topicName)
		}
	}

	return
}

func findPerOrgnaization(appName, org,
	topicName string) (topic *Topic, err error) {

	topic, err = Get(BLANK, appName, org, topicName)
	if err != nil {
		topic, err = Get(BLANK, BLANK, org, topicName)
	}

	return
}

func findGlobal(appName,
	topicName string) (topic *Topic, err error) {
	topic, err = Get(BLANK, appName, BLANK, topicName)
	if err != nil {
		topic, err = Get(BLANK, BLANK, BLANK, topicName)
	}

	return
}

func Find(user, appName, org,
	topicName string) (topic *Topic, err error) {

	topic, err = findPerUser(user, appName, org, topicName)
	if err != nil {
		topic, err = findPerOrgnaization(appName, org, topicName)
		if err != nil {
			topic, err = findGlobal(appName, topicName)
		}
	}

	return

}
