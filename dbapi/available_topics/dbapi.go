package available_topics

import (
	"log"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
)

func GetAppTopics(app_name, org string) []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	topics := []string{}

	err := db.Conn.GetCursor(session, db.AvailableTopicCollection, query).Distinct("ident", &topics)
	if err != nil {
		log.Println(err)
	}

	return topics
}

type ChotaTopic map[string]string

func GetAll(user, app_name, org string, channels []string) (map[string]ChotaTopic, error) {
	appTopics := GetAppTopics(app_name, org)
	topicMap := map[string]ChotaTopic{}

	for _, ident := range appTopics {
		ct := ChotaTopic{"topic": ident}
		for _, channel := range channels {
			ct[channel] = "true"
		}

		topicMap[ident] = ct
	}

	disabled := new(topics.Topic)

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

func Get(topic string) (found *db.AvailableTopic, err error) {
	found = new(db.AvailableTopic)

	err = db.Conn.GetOne(db.AvailableTopicCollection, utils.M{"ident": topic}, found)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func Insert(newTopic *db.AvailableTopic) string {
	return db.Conn.Insert(db.AvailableTopicCollection, newTopic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.AvailableTopicCollection, *doc)
}
