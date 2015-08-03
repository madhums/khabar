// Provides methods for creating and modifying preferenfes
package topics

import (
	"log"

	"github.com/bulletind/khabar/db"
	availableTopics "github.com/bulletind/khabar/dbapi/available_topics"
	"github.com/bulletind/khabar/utils"
)

// Insert adds a new user or organization preference
func Insert(topic *db.Topic) string {
	return db.Conn.Insert(db.TopicCollection, topic)
}

// InsertOrUpdateTopic adds a topic in `topics` collection if it doesn't exist
// Or it update the topic if it exists
//
// This is not a very effecient way of doing it but we are doing it at the cost
// of the data structure we want.
//
// One of the limitations is that mongo is not yet capable of upserting to array
// of documents. For this to happen in a simple way the data structure has to
// change from channels: []Channel to channels: map[string]interface
func InsertOrUpdateTopic(org, ident, channelName, attr string, val bool, user string) error {

	var channels []db.Channel
	var channel db.Channel
	var doc utils.M
	var spec utils.M

	found := new(db.Topic)
	query := utils.M{
		"org":   org,
		"user":  user, // empty for org
		"ident": ident,
	}

	// See if this setting already exists
	err := db.Conn.GetOne(
		db.TopicCollection,
		query,
		found,
	)

	// Fetch the global setting
	// Because if the above setting already doesn't exist, then set the other
	// attr to what is set in global

	global := new(db.Topic)
	q := utils.M{"org": "", "user": "", "ident": ident}
	e := db.Conn.GetOne(db.TopicCollection, q, &global)

	if e == nil {
		for _, ch := range global.Channels {
			if ch.Name == channelName {
				if attr == "Default" {
					channel.Locked = ch.Locked
				} else if attr == "Locked" {
					channel.Default = ch.Default
				}
			} else if err != nil {
				channels = append(channels, ch)
			}
		}
	}

	channel.Name = channelName

	if attr == "Default" {
		channel.Default = val
	} else if attr == "Locked" {
		channel.Locked = val
	} else {
		channel.Enabled = val
	}

	channels = append(channels, channel)

	// If it doesn't exist, insert and return

	if err != nil {
		topic := new(db.Topic)
		topic.PrepareSave()
		topic.Ident = ident
		topic.Organization = org
		topic.User = user
		topic.Channels = channels
		Insert(topic)
		return nil
	}

	// If it does exist, find the document in the array and modify it
	// Do one of the two depending on whether its present or not
	// Step 1. if its not present, add to channels array
	// Step 2. if its present, set the value

	query["channels.name"] = channelName
	err = db.Conn.GetOne(
		db.TopicCollection,
		query,
		found,
	)

	// Step 1. Add to set and return

	if err != nil {
		doc = utils.M{
			"$addToSet": utils.M{
				"channels": channel,
			},
		}
		delete(query, "channels.name")
		return updateTopics(query, doc)
	}

	// Step 2. Else set the value

	if attr == "Default" {
		spec = utils.M{
			"channels.$.default": val,
		}
	} else if attr == "Locked" {
		spec = utils.M{
			"channels.$.locked": val,
		}
	} else if attr == "Enabled" {
		spec = utils.M{
			"channels.$.enabled": val,
		}
	}

	doc = utils.M{
		"$set": spec,
	}

	return updateTopics(query, doc)
}

// updateTopics updates user or organization preferences
func updateTopics(query, doc utils.M) error {
	return db.Conn.Update(db.TopicCollection, query, doc)
}

// Initialize creates a list of "default - non-enabled" preferences for user or org
func Initialize(user, org string) error {
	if user == db.BLANK {
		org = db.BLANK
	}

	disabled := GetAllDefault(org)
	preferences := []interface{}{}

	for _, topic := range disabled {
		preference := db.Topic{
			User:         user,
			Organization: org,
			Ident:        topic.Ident,
			Channels:     topic.Channels,
		}

		preference.PrepareSave()
		preferences = append(preferences, preference)
	}

	if len(preferences) > 0 {
		err, _ := db.Conn.InsertMulti(db.TopicCollection, preferences...)
		return err
	}

	return nil
}

// GetAllDefault returns all the "default - non-enabled" preferences for the organization
func GetAllDefault(org string) []db.Topic {
	session := db.Conn.Session.Copy()
	defer session.Close()

	result := []db.Topic{}

	db.Conn.Get(session, db.TopicCollection, utils.M{
		"org":              org,
		"user":             db.BLANK,
		"channels.default": false,
	}).All(&result)

	return result
}

// ChannelAllowed checks if the requested channel is allowed by the user for sending
// out the notification
func ChannelAllowed(user, org, app_name, ident, channelName string) bool {

	var available = []string{"email", "web", "push"}
	var preference map[string]availableTopics.ChotaTopic

	appTopics, err := availableTopics.GetAppTopics(app_name, org)
	channels := []string{}
	for _, idnt := range available {
		channels = append(channels, idnt)
	}
	preference, err = availableTopics.GetUserPreferences(user, org, appTopics, &channels)

	if err != nil {
		log.Println(err)
		return false
	}

	if _, ok := preference[ident]; !ok {
		return false
	}

	if _, ok := preference[ident][channelName]; !ok {
		return false
	}

	return preference[ident][channelName].Value
}

func DisableUserChannel(orgs, topics []string, user, channelName string) {
	session := db.Conn.Session.Copy()
	defer session.Close()

	utils.RemoveDuplicates(&orgs)
	utils.RemoveDuplicates(&topics)

	db.Conn.Update(
		db.TopicCollection, utils.M{"user": user},
		utils.M{"$addToSet": utils.M{"channels": channelName}},
	)

	disabled := []interface{}{}

	for _, org := range orgs {
		disabledTopics := []string{}
		db.Conn.GetCursor(session, db.TopicCollection, utils.M{"user": user, "org": org}).Distinct("ident", &disabledTopics)

		for _, name := range topics {
			if !db.InArray(name, disabledTopics) {
				topic := db.Topic{
					User:         user,
					Organization: org,
					Ident:        name,
					Channels: []db.Channel{
						db.Channel{Name: channelName, Enabled: false},
					},
				}

				topic.PrepareSave()
				disabled = append(disabled, topic)
			}
		}
	}

	if len(disabled) > 0 {
		db.Conn.InsertMulti(db.TopicCollection, disabled...)
	}
}

func Get(user, org, topicName string) (topic *db.Topic, err error) {

	topic = new(db.Topic)

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

// DeleteTopic deletes the ident from `topics_available` collection
func DeleteTopic(ident string) error {
	err := db.Conn.Delete(db.TopicCollection, utils.M{"ident": ident})
	if err != nil {
		return err
	}
	err = db.Conn.Delete(db.AvailableTopicCollection, utils.M{"ident": ident})
	return err
}

// AddChannel enables the channel for that particular ident for sending notifications
func AddChannel(ident, channelName, user, organization string) error {
	query := utils.M{
		"org":   organization,
		"user":  user,
		"ident": ident,
	}

	spec := utils.M{
		"$addToSet": utils.M{
			"channels": db.Channel{Name: channelName, Enabled: true},
		},
	}

	result := utils.M{}

	_, err := db.Conn.FindAndUpdate(db.TopicCollection, query, spec, &result)
	return err
}

// RemoveChannel disables the channel for that particular ident for sending
// notifications
func RemoveChannel(ident, channelName, user, organization string) error {
	query := utils.M{
		"org":           organization,
		"user":          user,
		"ident":         ident,
		"channels.name": channelName,
	}

	found := new(db.Topic)

	err := db.Conn.GetOne(
		db.TopicCollection,
		query,
		found,
	)

	err = db.Conn.Update(
		db.TopicCollection,
		query,
		utils.M{
			"$set": utils.M{
				"channels.$.enabled": false,
			},
		},
	)

	return err
}
