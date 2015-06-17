package available_topics

import (
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

func GetAppTopics(app_name, org string) *[]string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	topics := []string{}

	var topic struct {
		Ident string `bson:"ident"`
	}

	iter := db.Conn.GetCursor(
		session, db.AvailableTopicCollection, query,
	).Select(utils.M{"ident": 1}).Sort("ident").Iter()

	for iter.Next(&topic) {
		topics = append(topics, topic.Ident)
	}

	return &topics
}

func GetOrgTopics(org string, appTopics, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for org level

	topicMap := map[string]ChotaTopic{}

	for _, ident := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range *channels {
			ct[channel] = &TopicDetail{Locked: false, Value: trueState}
		}

		topicMap[ident] = ct
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
				topicMap[disabled.Ident][blocked].Value = disabledState
			}
		}
	}

	return topicMap, nil
}

func ApplyLockes(org string, appTopics *[]string, topicMap map[string]ChotaTopic) {
	enabled := locks.GetAll(org)
	for _, pref := range enabled {
		if _, ok := topicMap[pref.Topic]; ok {
			for _, blocked := range pref.Channels {

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

func GetUserTopics(user, org string, appTopics, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for user level

	topicMap := map[string]ChotaTopic{}

	for _, ident := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range *channels {
			ct[channel] = &TopicDetail{Locked: false, Value: trueState}
		}

		topicMap[ident] = ct
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
				topicMap[disabled.Ident][blocked].Value = falseState
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
				topicMap[disabled.Ident][blocked].Value = disabledState
			}
		}
	}

	ApplyLockes(org, appTopics, topicMap)

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
