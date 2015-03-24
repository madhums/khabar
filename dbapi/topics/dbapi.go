package topics

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

const BLANK = ""

func Update(dbConn *db.MConn, user string, appName string, org string, topicName string, doc *utils.M) error {

	return dbConn.Update(db.TopicCollection,
		utils.M{"app_name": appName,
			"org":   org,
			"user":  user,
			"ident": topicName,
		},
		utils.M{
			"$set": *doc,
		})
}

func Insert(dbConn *db.MConn, topic *Topic) string {
	return dbConn.Insert(db.TopicCollection, topic)
}

func Delete(dbConn *db.MConn, doc *utils.M) error {
	return dbConn.Delete(db.TopicCollection, *doc)
}

func Get(dbConn *db.MConn, user, appName, org, topicName string) (topic *Topic, err error) {

	topic = new(Topic)

	err = dbConn.GetOne(
		db.TopicCollection,
		utils.M{
			"app_name": appName,
			"org":      org,
			"user":     user,
			"ident":    topicName,
		},
		topic,
	)

	if err != nil {
		return nil, err
	}

	return
}

func GetAll(dbConn *db.MConn, user, appName, org string) (*[]Topic, error) {
	var query utils.M = make(utils.M)

	var result []Topic

	if len(user) > 0 {
		query["user"] = user
	}

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	session := dbConn.Session.Copy()
	defer session.Close()

	err := dbConn.GetCursor(session, db.TopicCollection, query).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func findPerUser(dbConn *db.MConn, user, appName, org, topicName string) (topic *Topic, err error) {

	topic, err = Get(dbConn, user, appName, org, topicName)
	if err != nil {
		topic, err = Get(dbConn, user, appName, BLANK, topicName)
		if err != nil {
			topic, err = Get(dbConn, user, BLANK, org, topicName)
		}
	}

	return
}

func findPerOrgnaization(dbConn *db.MConn, appName, org, topicName string) (topic *Topic, err error) {

	topic, err = Get(dbConn, BLANK, appName, org, topicName)
	if err != nil {
		topic, err = Get(dbConn, BLANK, BLANK, org, topicName)
	}

	return
}

func findGlobal(dbConn *db.MConn, appName, topicName string) (topic *Topic, err error) {
	topic, err = Get(dbConn, BLANK, appName, BLANK, topicName)
	if err != nil {
		topic, err = Get(dbConn, BLANK, BLANK, BLANK, topicName)
	}

	return
}

func Find(dbConn *db.MConn, user, appName, org, topicName string) (topic *Topic, err error) {

	topic, err = findPerUser(dbConn, user, appName, org, topicName)
	if err != nil {
		topic, err = findPerOrgnaization(dbConn, appName, org, topicName)
		if err != nil {
			topic, err = findGlobal(dbConn, appName, topicName)
		}
	}

	return

}
