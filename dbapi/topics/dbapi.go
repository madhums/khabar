package topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

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
			utils.M{"user": db.BLANK, "org": org},
			utils.M{"user": db.BLANK, "org": db.BLANK},
			utils.M{"user": user, "org": db.BLANK},
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

func findPerUser(user, org, topicName string) (topic *Topic, err error) {

	topic, err = Get(user, org, topicName)
	if err != nil {
		topic, err = Get(user, db.BLANK, topicName)
	}

	return
}

func findPerOrgnaization(org, topicName string) (topic *Topic, err error) {
	return Get(db.BLANK, org, topicName)
}

func findGlobal(topicName string) (topic *Topic, err error) {
	return Get(db.BLANK, db.BLANK, topicName)
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

func DeleteTopic(ident string) error {
	err := db.Conn.Delete(db.TopicCollection, utils.M{"ident": ident})
	if err != nil {
		return err
	}
	err = db.Conn.Delete(db.AvailableTopicCollection, utils.M{"ident": ident})
	return err
}
