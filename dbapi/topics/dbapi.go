package topics

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

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

func Get(dbConn *db.MConn, user string, appName string, org string, topicName string) (topic *Topic, err error) {

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

func GetAll(dbConn *db.MConn, user string, appName string, org string) (*[]Topic, error) {
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

func findPerUser(dbConn *db.MConn, user string, appName string, org string, topicName string) (topic *Topic, err error) {

	topic, err = Get(dbConn, user, appName, org, topicName)

	if err == nil {
		return
	}

	topic, err = Get(dbConn, user, appName, "", topicName)

	if err == nil {
		return
	}

	topic, err = Get(dbConn, user, "", org, topicName)

	if err == nil {
		return
	}
	return
}

func findPerOrgnaization(dbConn *db.MConn, appName string, org string, topicName string) (topic *Topic, err error) {

	topic, err = Get(dbConn, "", appName, org, topicName)

	if err == nil {
		return
	}

	topic, err = Get(dbConn, "", "", org, topicName)

	if err == nil {
		return
	}

	return

}

func findGlobal(dbConn *db.MConn, topicName string) (topic *Topic, err error) {

	topic, err = Get(dbConn, "", "", "", topicName)

	if err == nil {
		return
	}

	return

}

func Find(dbConn *db.MConn, user string, appName string, org string, topicName string) (topic *Topic, err error) {

	topic, err = findPerUser(dbConn, user, appName, org, topicName)

	if err == nil {
		return
	}

	topic, err = findPerOrgnaization(dbConn, appName, org, topicName)

	if err == nil {
		return
	}

	topic, err = findGlobal(dbConn, topicName)

	if err == nil {
		return
	}

	return

}
