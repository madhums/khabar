package topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/available_topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
)

const BLANK = ""

func Update(user, org, topicName string, doc *utils.M) error {

	return db.Conn.Update(db.TopicCollection,
		utils.M{
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

func ChannelAllowed(user, org, topicName, channel string) bool {
	return db.Conn.Count(db.TopicCollection, utils.M{
		"$or": []utils.M{
			utils.M{"user": BLANK, "org": org},
			utils.M{"user": BLANK, "org": BLANK},
			utils.M{"user": user, "org": BLANK},
			utils.M{"user": user, "org": org},
		},
		"ident":    topicName,
		"channels": channel,
	}) == 0
}

func Get(user, org, topicName string) (topic *Topic, err error) {

	topic = new(Topic)

	err = db.Conn.GetOne(
		db.TopicCollection,
		utils.M{
			"org":   org,
			"user":  user,
			"ident": topicName,
		},
		topic,
	)

	if err != nil {
		return nil, err
	}

	return
}

func GetAll(user, app_name, org string) (*mgo.Iter, error) {
	appTopics := available_topics.GetAppTopics(app_name, org)

	var query utils.M = make(utils.M)

	query["ident"] = utils.M{"$in": appTopics}
	query["user"] = user
	query["org"] = org

	session := db.Conn.Session.Copy()
	defer session.Close()

	iter := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return iter, nil
}

func findPerUser(user, org, topicName string) (topic *Topic, err error) {

	topic, err = Get(user, org, topicName)
	if err != nil {
		topic, err = Get(user, BLANK, topicName)
	}

	return
}

func findPerOrgnaization(org, topicName string) (topic *Topic, err error) {
	return Get(BLANK, org, topicName)
}

func findGlobal(topicName string) (topic *Topic, err error) {
	return Get(BLANK, BLANK, topicName)
}

func Find(user, org, topicName string) (topic *Topic, err error) {

	topic, err = findPerUser(user, org, topicName)
	if err != nil {
		topic, err = findPerOrgnaization(org, topicName)
		if err != nil {
			topic, err = findGlobal(topicName)
		}
	}

	return
}
