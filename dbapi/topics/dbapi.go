package topics

import (
	"log"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
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

func getAppTopics(app_name, org string) []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	topics := []string{}

	err := db.Conn.GetCursor(session, db.TopicsAvailable, query).Distinct("ident", &topics)
	if err != nil {
		log.Println(err)
	}

	return topics
}

type ChotaTopic map[string]string

func GetAll(user, app_name, org string, channels []string) (map[string]ChotaTopic, error) {
	appTopics := getAppTopics(app_name, org)
	topicMap := map[string]ChotaTopic{}

	for _, ident := range appTopics {
		ct := ChotaTopic{"topic": ident}
		for _, channel := range channels {
			ct[channel] = "true"
		}

		topicMap[ident] = ct
	}

	disabled := new(Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{
		"ident": utils.M{"$in": appTopics},
		"user":  user,
		"org":   org,
	}

	userBlacklisted := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for userBlacklisted.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked] = "false"
			}
		}
	}

	delete(query, "user")

	orgBlacklisted := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for orgBlacklisted.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked] = "disabled"
			}
		}
	}

	return topicMap, nil
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

func DeleteTopic(ident string) {
	db.Conn.Delete(db.TopicCollection, utils.M{"ident": ident})
	db.Conn.Delete(db.TopicsAvailable, utils.M{"ident": ident})
}
