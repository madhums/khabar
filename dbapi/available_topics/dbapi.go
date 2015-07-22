package available_topics

import (
	"encoding/json"
	"fmt"

	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/locks"
	"gopkg.in/bulletind/khabar.v1/utils"
)

const trueState = "true"
const falseState = "false"
const disabledState = "disabled"

type TopicDetail struct {
	Locked bool   `json:"locked"`
	Value  string `json:"value"`
}

type ChotaTopic map[string]*TopicDetail

func GetAllTopics() []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	topics := []string{}

	db.Conn.GetCursor(
		session, db.AvailableTopicCollection, utils.M{},
	).Distinct("ident", &topics)

	return topics
}

func GetAppTopics(app_name, org string) *[]db.AvailableTopic {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	var topics []db.AvailableTopic

	iter := db.Conn.GetCursor(
		session, db.AvailableTopicCollection, query,
	).Select(utils.M{"ident": 1, "channels": 1}).Sort("ident").Iter()

	fmt.Println("====================")
	o, err := json.Marshal(topics)
	if err != nil {
		fmt.Println("Error marshaling JSON")
	}
	fmt.Println(string(o))
	fmt.Println("=================")

	return &topics
}

func GetOrgTopics(org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for org level

	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: trueState}
		}

		topicMap[availableTopic.Ident] = ct
	}

	disabled := new(db.Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{
		"ident": utils.M{"$in": appTopics},
		"user":  db.BLANK,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked].Value = falseState
			}
		}
	}

	//Find Globally Disabled Topics
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				delete(topicMap[disabled.Ident], blocked)
			}
		}
	}

	return topicMap, nil
}

func ApplyLocks(org string, topicMap map[string]ChotaTopic) {
	enabled := locks.GetAll(org)
	for _, pref := range enabled {
		if _, ok := topicMap[pref.Topic]; ok {
			for _, blocked := range pref.Channels {

				if _, ok := topicMap[pref.Topic][blocked]; !ok {
					continue
				}

				if topicMap[pref.Topic][blocked].Value == disabledState {
					continue
				}

				topicMap[pref.Topic][blocked].Locked = true

				if pref.Enabled {
					topicMap[pref.Topic][blocked].Value = trueState
				} else {
					topicMap[pref.Topic][blocked].Value = falseState
				}
			}
		}
	}
}

func GetUserTopics(user, org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for user level

	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: falseState}
		}

		topicMap[availableTopic.Ident] = ct
	}

	disabled := new(db.Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{
		"ident": utils.M{"$in": appTopics},
		"user":  user,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked].Value = trueState
			}
		}
	}

	//Find all Topics that have been blocked by the Organization
	query["user"] = db.BLANK

	pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass2.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked].Value = disabledState
			}
		}
	}

	//Find Globally Disabled Topics
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				delete(topicMap[disabled.Ident], blocked)
			}
		}
	}

	ApplyLocks(org, topicMap)

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
