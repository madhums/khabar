package topics

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

func Update(user, org, ident string, doc *utils.M) error {

	return db.Conn.Update(db.TopicCollection,
		utils.M{
			"org":   org,
			"user":  user,
			"ident": ident,
		},
		utils.M{
			"$set": *doc,
		})
}

func Insert(topic *db.Topic) string {
	return db.Conn.Insert(db.TopicCollection, topic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.TopicCollection, *doc)
}

/**
 * Insert a topic in `topics` collection if it doesn't exist
 * Or Update the topic
 *
 * - 	This is not a very effecient way of doing it but we are doing it at
 * 		the cost of the data structure we want
 * - 	One of the limitations is that mongo is not yet capable of
 * 		upserting to array of documents. For this to happen in a simple way the data
 * 		structure has to change OR we use multiple collections
 */

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
		return AddOrgOrUserChannel(query, doc)
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

	return AddOrgOrUserChannel(query, doc)
}

func AddOrgOrUserChannel(query, doc utils.M) error {
	return db.Conn.Update(db.TopicCollection, query, doc)
}

func GetChannelProperty(channels []db.Channel, channelName, attr string) bool {
	var val bool
	for _, channel := range channels {
		if channel.Name == channelName {
			if attr == "Default" {
				val = channel.Default
			} else if attr == "Locked" {
				val = channel.Locked
			} else if attr == "Enabled" {
				val = channel.Enabled
			}
		}
	}
	return val
}

func Initialize(user, org string) error {
	if user == db.BLANK {
		org = db.BLANK
	}

	disabled := GetAllDisabled(org)
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

func GetAllDisabled(org string) []db.Topic {
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

func ChannelAllowed(user, org, topicName, channel string) bool {
	// TODO: Populate default values here.
	defUser := db.Conn.Count(db.TopicCollection, utils.M{
		"$or": []utils.M{
			utils.M{"user": user, "org": db.BLANK},
			utils.M{"user": user, "org": org},
		},
		"ident":    topicName,
		"channels": channel,
	}) == 0

	defOrg := db.Conn.Count(db.TopicCollection, utils.M{
		"$or": []utils.M{
			utils.M{"user": db.BLANK, "org": org},
			utils.M{"user": db.BLANK, "org": db.BLANK},
		},
		"ident":    topicName,
		"channels": channel,
	}) == 0

	lockEntry := new(db.Locks)

	err := db.Conn.GetOne(db.LocksCollection,
		utils.M{"org": org, "ident": topicName, "channels": channel},
		lockEntry)

	if err != nil {
		return defOrg && defUser
	} else {
		return defOrg && lockEntry.Enabled
	}
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

func findPerUser(user, org, topicName string) (topic *db.Topic, err error) {

	topic, err = Get(user, org, topicName)
	if err != nil {
		topic, err = Get(user, db.BLANK, topicName)
	}

	return
}

func findPerOrgnaization(org, topicName string) (topic *db.Topic, err error) {
	return Get(db.BLANK, org, topicName)
}

func findGlobal(topicName string) (topic *db.Topic, err error) {
	return Get(db.BLANK, db.BLANK, topicName)
}

func Find(user, org, topicName string) (topic *db.Topic, err error) {

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
