package topics

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/defaults"
	"gopkg.in/bulletind/khabar.v1/utils"
)

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

func Insert(topic *db.Topic) string {
	return db.Conn.Insert(db.TopicCollection, topic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.TopicCollection, *doc)
}

/**
 * Insert a topic in `topics` collection if it doesn't exist
 * Or Update the topic
 * Used only in the org level mostly to set default
 */

func InsertOrUpdateTopic(org, ident string, channelName string) error {

	found := new(db.Topic)
	query := utils.M{
		"org":           org,
		"user":          "",
		"ident":         ident,
		"channels.name": channelName,
	}
	channels := []db.Channel{
		db.Channel{Name: channelName, Default: true},
	}

	err := db.Conn.GetOne(
		db.TopicCollection,
		query,
		found,
	)

	// If it doesn't exist, insert and return
	if err != nil {
		topic := new(db.Topic)
		topic.PrepareSave()
		// topic.ToggleValue() // default `value` is false, so toggle it
		topic.Ident = ident
		topic.Organization = org
		topic.Channels = channels
		Insert(topic)
		return nil
	}

	// Update Default attribute (toggle it)

	err = db.Conn.Update(
		db.TopicCollection,
		query,
		utils.M{
			"$set": utils.M{
				"channels.default": !found.Channels[0].Default,
			},
		},
	)

	return err
}

func Initialize(user, org string) error {
	orgArg := org
	if user == db.BLANK {
		orgArg = db.BLANK
	}

	disabled := defaults.GetAllDisabled(orgArg)

	preferences := []interface{}{}

	for _, entry := range disabled {
		preference := db.Topic{
			User:         user,
			Organization: org,
			Ident:        entry.Topic,
			Channels:     entry.Channels,
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

func DisableUserChannel(orgs, topics []string, user, channel string) {
	session := db.Conn.Session.Copy()
	defer session.Close()

	utils.RemoveDuplicates(&orgs)
	utils.RemoveDuplicates(&topics)

	db.Conn.Update(
		db.TopicCollection, utils.M{"user": user},
		utils.M{"$addToSet": utils.M{"channels": channel}},
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
					Channels:     []string{channel},
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

func AddChannel(ident, channel, user, organization string) error {
	query := utils.M{
		"org":   organization,
		"user":  user,
		"ident": ident,
	}

	spec := utils.M{"$addToSet": utils.M{"channels": channel}}

	result := utils.M{}

	_, err := db.Conn.FindAndUpdate(db.TopicCollection, query, spec, &result)
	return err
}

func RemoveChannel(ident, channel, user, organization string) error {
	query := utils.M{
		"org":   organization,
		"user":  user,
		"ident": ident,
	}

	spec := utils.M{"$pull": utils.M{"channels": channel}}

	result := utils.M{}

	_, err := db.Conn.FindAndUpdate(db.TopicCollection, query, spec, &result)
	return err
}
