package available_topics

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/locks"
	"gopkg.in/bulletind/khabar.v1/utils"
)

const trueState = true
const falseState = false
const disabledState = false

type TopicDetail struct {
	Locked  bool `json:"locked"`
	Value   bool `json:"value"`
	Default bool `json:"default"`
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

func GetAppTopics(app_name, org string) (*[]db.AvailableTopic, error) {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	var topics []db.AvailableTopic

	iter := db.Conn.GetCursor(
		session, db.AvailableTopicCollection, query,
	).Select(utils.M{"ident": 1, "channels": 1}).Sort("ident").Iter()

	err := iter.All(&topics)
	// TODO: handle this error

	return &topics, err
}

func GetOrgTopics(org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for org level
	var availableTopics []string

	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: trueState}
		}

		topicMap[availableTopic.Ident] = ct
	}

	topic := new(db.Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	for _, topic := range *appTopics {
		availableTopics = append(availableTopics, topic.Ident)
	}

	query := utils.M{
		"ident": utils.M{"$in": availableTopics},
		"user":  db.BLANK,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(topic) {

		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				topicMap[topic.Ident][channel].Default = topic.Value
				// topicMap[topic.Ident][channel].Locked = topic.Value
			}
		}
	}

	//Find Globally topic Topics
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				delete(topicMap[topic.Ident], channel)
			}
		}
	}

	ApplyLocks(org, topicMap)

	return topicMap, nil
}

func ApplyLocks(org string, topicMap map[string]ChotaTopic) {
	enabled := locks.GetAll(org)

	for _, pref := range enabled {
		if _, ok := topicMap[pref.Topic]; ok {
			for _, channel := range pref.Channels {

				if _, ok := topicMap[pref.Topic][channel]; !ok {
					continue
				}

				// if topicMap[pref.Topic][channel].Value == disabledState {
				// 	continue
				// }

				topicMap[pref.Topic][channel].Locked = true

				// if pref.Enabled {
				// 	topicMap[pref.Topic][channel].Value = trueState
				// } else {
				// 	topicMap[pref.Topic][channel].Value = falseState
				// }
			}
		}
	}

}

// TODO: There's quite some repetetive code, break it down

func GetUserTopics(user, org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {

	// This is a bit clumsy, make it better later
	// TODO: Explain the usage
	topicmap2 := map[string]string{}

	var availableTopics []string
	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: falseState}
		}

		topicMap[availableTopic.Ident] = ct
	}

	topic := new(db.Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	for _, topic := range *appTopics {
		availableTopics = append(availableTopics, topic.Ident)
	}

	// Step 1
	// Add defaults for user level
	query := utils.M{
		"ident": utils.M{"$in": availableTopics},
		"user":  user,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				topicmap2[topic.Ident] = channel
				topicMap[topic.Ident][channel].Value = trueState
			}
		}
	}

	// Step 2
	// Find all Topics that have been defaulted by the Organization
	query["user"] = db.BLANK

	pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass2.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				// This should be overridden by `topic.Value`.
				// But before doing that, we must make sure the user hasn't ever set
				// this field.

				topicMap[topic.Ident][channel].Value = disabledState
				topicMap[topic.Ident][channel].Default = topic.Value

				for _, value := range topicmap2 {
					if value == channel {
						topicMap[topic.Ident][channel].Value = topic.Value
					}
				}
			}
		}
	}

	// Step 3
	// Find Globally topic Topics
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(topic) {

		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				delete(topicMap[topic.Ident], channel)
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
