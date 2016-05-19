// Provides methods for fetching preferences and modifying availbale topics
package available_topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

type TopicDetail struct {
	Locked  bool `json:"locked"`
	Value   bool `json:"value"`
	Default bool `json:"default"`
}

type ChotaTopic map[string]*TopicDetail

// GetAllTopics returns all the available topics
func GetAllTopics() []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	topics := []string{}

	db.Conn.GetCursor(
		session, db.AvailableTopicCollection, utils.M{},
	).Distinct("ident", &topics)

	return topics
}

// GetAppTopics returns all the available topics for the particular app
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

// GetOrgPreferences lists all the organization preferences for the ident and channel
func GetOrgPreferences(org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {
	// Add defaults for org level
	var availableTopics []string

	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: true}
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
				// topicMap[topic.Ident][channel.Name].Value = channel.Enabled
				topicMap[topic.Ident][channel.Name].Locked = channel.Locked
			}
		}
	}

	return topicMap, nil
}

// GetUserPreferences lists all the user preferences
func GetUserPreferences(user, org string, appTopics *[]db.AvailableTopic, channels *[]string) (map[string]ChotaTopic, error) {

	// Remember what the original org and global settings for ident x channel
	orgSetting := map[string]ChotaTopic{}

	var availableTopics []string
	topicMap := map[string]ChotaTopic{}

	for _, availableTopic := range *appTopics {
		ct := ChotaTopic{}
		for _, channel := range availableTopic.Channels {
			ct[channel] = &TopicDetail{Locked: false, Value: false}
		}

		orgSetting[availableTopic.Ident] = ct
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
				topicMap[topic.Ident][channel.Name].Value = channel.Default
			}
		}
	}

	// Step 2
	// Use org settings

	query["user"] = db.BLANK
	query["org"] = org

	pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass2.Next(topic) {
		if _, ok := topicMap[topic.Ident]; ok {
			for _, channel := range topic.Channels {

				if _, ok = topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				// Remember org setting
				orgSetting[topic.Ident][channel.Name].Default = channel.Default
				orgSetting[topic.Ident][channel.Name].Locked = channel.Locked

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

				if _, ok = topicMap[topic.Ident][channel.Name]; !ok {
					continue
				}

				// This is what the user has set
				topicMap[topic.Ident][channel.Name].Value = channel.Enabled

				// No need to override the lock if the org setting doesn't exist for this channel
				if _, ok = orgSetting[topic.Ident][channel.Name]; !ok {
					continue
				}

				if orgSetting[topic.Ident][channel.Name].Default && orgSetting[topic.Ident][channel.Name].Locked {
					topicMap[topic.Ident][channel.Name].Value = true
				} else if !orgSetting[topic.Ident][channel.Name].Default && orgSetting[topic.Ident][channel.Name].Locked {
					topicMap[topic.Ident][channel.Name].Value = false
				}
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

// Insert creates a new available topic
func Insert(newTopic *db.AvailableTopic) string {
	return db.Conn.Insert(db.AvailableTopicCollection, newTopic)
}

// Delete removes an available topic
func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.AvailableTopicCollection, *doc)
}
