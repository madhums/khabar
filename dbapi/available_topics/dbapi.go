package available_topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
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

func GetAllLocked(org string) []db.Topic {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Topic{}

	db.Conn.Get(session, db.TopicCollection, utils.M{
		"org":             org,
		"user":            "",
		"channels.locked": true,
	}).All(&result)

	return result
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
				topicMap[topic.Ident][channel.Name].Default = channel.Default
				// topicMap[topic.Ident][channel].Locked = topic.Value
			}
		}
	}

	// Find Globally disabled Topics
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				delete(topicMap[topic.Ident], channel.Name)
			}
		}
	}

	ApplyLocks(org, topicMap)

	return topicMap, nil
}

func ApplyLocks(org string, topicMap map[string]ChotaTopic) {
	locked := GetAllLocked(org)

	for _, topic := range locked {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				if _, ok := topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				topicMap[topic.Ident][channel.Name].Locked = channel.Locked
			}
		}
	}
}

func GetUserTopics(user, org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {

	// We are trying to remember what the original user setting was for ident x channel
	userSetting := make(map[string][]db.Channel)

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
	// Add user preferences
	query := utils.M{
		"ident": utils.M{"$in": availableTopics},
		"user":  user,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				userSetting[topic.Ident] = topic.Channels

				// These is what the user has set
				topicMap[topic.Ident][channel.Name].Value = trueState
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
				// Set the default
				topicMap[topic.Ident][channel.Name].Default = channel.Default
			}
		}
	}

	// Step 3
	// Remove globablly disabled topic/channels
	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(topic) {

		// Override it with the global setting
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				delete(topicMap[topic.Ident], channel.Name)
			}
		}
	}

	// After all the overrides apply locks
	ApplyLocks(org, topicMap)

	// After the locks have been applied, make sure that the defaults are
	// applied properly

	for idnt, values := range topicMap {
		for ch, _ := range values {

			if topicMap[idnt][ch].Default && topicMap[idnt][ch].Locked {
				topicMap[idnt][ch].Value = true
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
