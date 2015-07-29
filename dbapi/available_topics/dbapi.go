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
		"org":   db.BLANK,
	}

	// Step 1
	// Use global settings

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				if _, ok = topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				topicMap[topic.Ident][channel.Name].Default = channel.Default
				topicMap[topic.Ident][channel.Name].Locked = channel.Locked
			}
		}
	}

	// Step 2
	// Override it with organization settings

	query["org"] = org

	pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass2.Next(topic) {

		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				topicMap[topic.Ident][channel.Name].Default = channel.Default
				// topicMap[topic.Ident][channel].Locked = topic.Value
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
	// userSetting := make(map[string][]db.Channel)

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

	query := utils.M{
		"ident": utils.M{"$in": availableTopics},
	}

	// Step 1
	// Add global settings

	query["user"] = db.BLANK
	query["org"] = db.BLANK

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(topic) {

		// Override it with the global setting
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {
				if _, ok = topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				topicMap[topic.Ident][channel.Name].Default = channel.Default
				topicMap[topic.Ident][channel.Name].Locked = channel.Locked
				topicMap[topic.Ident][channel.Name].Value = channel.Enabled
			}
		}
	}

	// Step 2
	// Use org settings

	query["user"] = db.BLANK

	pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass2.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				if _, ok = topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				// Set the default from org
				topicMap[topic.Ident][channel.Name].Default = channel.Default
				topicMap[topic.Ident][channel.Name].Locked = channel.Locked
				topicMap[topic.Ident][channel.Name].Value = channel.Default
			}
		}
	}

	// Step 3
	// Add user settings

	query["user"] = user
	query["org"] = org

	pass3 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass3.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				// userSetting[topic.Ident] = topic.Channels

				// This is what the user has set
				topicMap[topic.Ident][channel.Name].Value = channel.Enabled

				if topicMap[topic.Ident][channel.Name].Default && topicMap[topic.Ident][channel.Name].Locked {
					topicMap[topic.Ident][channel.Name].Value = true
				}
			}
		}
	}

	// After all the overrides apply locks
	// ApplyLocks(org, topicMap)

	// After the locks have been applied, make sure that the defaults are
	// applied properly

	// for idnt, values := range topicMap {
	// 	for ch, _ := range values {

	// 		if topicMap[idnt][ch].Default && topicMap[idnt][ch].Locked {
	// 			topicMap[idnt][ch].Value = true
	// 		}
	// 	}
	// }

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
