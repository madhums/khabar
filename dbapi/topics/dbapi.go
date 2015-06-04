package topics

import (
	"errors"

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

func Initialize(user, org string) {
	disabled := defaults.GetAllDisabled(user, org)

	preferences := []interface{}{}

	for _, entry := range disabled {
		preference := db.Topic{
			User:         user,
			Organization: org,
			Ident:        entry.Topic,
		}
	}

}

func ChannelAllowed(user, org, topicName, channel string) bool {
	return db.Conn.Count(db.TopicCollection, utils.M{
		"$or": []utils.M{
			utils.M{"user": db.BLANK, "org": org},
			utils.M{"user": db.BLANK, "org": db.BLANK},
			utils.M{"user": user, "org": db.BLANK},
			utils.M{"user": user, "org": org},
		},
		"ident":    topicName,
		"channels": channel,
	}) == 0
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
	if organization == db.BLANK && user == db.BLANK {
		return errors.New("Atleast one of the user or org must be present.")
	}

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

func RemoveChannel(ident, channel, user, organization string) error {
	if organization == db.BLANK && user == db.BLANK {
		return errors.New("Atleast one of the user or org must be present.")
	}

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
