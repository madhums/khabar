package topics

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
	"gopkg.in/simversity/gottp.v2"
)

func Update(dbConn *db.MConn, user string, appName string, org string, topicName string, doc *utils.M) error {

	return dbConn.Update(TopicCollection,
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
	return dbConn.Insert(TopicCollection, topic)
}

func Delete(dbConn *db.MConn, doc *utils.M) error {
	return dbConn.Delete(TopicCollection, *doc)
}

func Get(dbConn *db.MConn, user string, appName string, org string, topicName string) *Topic {
	topic := new(Topic)
	if dbConn.GetOne(TopicCollection, utils.M{"app_name": appName,
		"org": org, "user": user, "ident": topicName}, topic) != nil {
		return nil
	}
	return topic
}

func GetAll(dbConn *db.MConn, paginator *gottp.Paginator, user string, appName string, org string, ident string) *[]Topic {
	var query utils.M = make(utils.M)
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
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

	if len(ident) > 0 {
		query["ident"] = ident
	}

	var limit int
	var skip int
	var limitExists bool
	var skipExists bool

	limit, limitExists = query["limit"].(int)
	skip, skipExists = query["skip"].(int)

	delete(query, "limit")
	delete(query, "skip")

	if !limitExists {
		limit = 30
	}
	if !skipExists {
		skip = 0
	}

	delete(query, "limit")
	delete(query, "skip")

	dbConn.GetCursor(TopicCollection, query).Skip(skip).Limit(limit).All(&result)

	return &result
}

func findPerUser(dbConn *db.MConn, user string, appName string, org string, topicName string) *Topic {
	var err error
	var topic *Topic

	topic = Get(dbConn, user, appName, org, topicName)

	if topic != nil {
		return topic
	}

	err = dbConn.GetOne(TopicCollection, utils.M{
		"user":     user,
		"app_name": appName,
		"org":      "",
		"ident":    topicName,
	}, topic)

	if err == nil {
		return topic
	}

	err = dbConn.GetOne(TopicCollection, utils.M{
		"user":     user,
		"app_name": "",
		"org":      org,
		"ident":    topicName,
	}, topic)

	if err == nil {
		return topic
	}
	return nil
}

func findPerOrgnaization(dbConn *db.MConn, appName string, org string, topicName string) *Topic {
	var err error
	topic := new(Topic)
	err = dbConn.GetOne(TopicCollection, utils.M{
		"user":     "",
		"app_name": appName,
		"org":      org,
		"ident":    topicName,
	}, topic)

	if err == nil {
		return topic
	}

	err = dbConn.GetOne(TopicCollection, utils.M{
		"user":     "",
		"app_name": "",
		"org":      org,
		"ident":    topicName,
	}, topic)

	if err == nil {
		return topic
	}

	return nil

}

func findGlobal(dbConn *db.MConn, topicName string) *Topic {
	var err error
	topic := new(Topic)
	err = dbConn.GetOne(TopicCollection, utils.M{
		"user":     "",
		"app_name": "",
		"org":      "",
		"ident":    topicName,
	}, topic)

	if err == nil {
		return topic
	}

	return nil

}

func Find(dbConn *db.MConn, user string, appName string, org string, topicName string) *Topic {
	var topic *Topic

	topic = findPerUser(dbConn, user, appName, org, topicName)

	if topic != nil {
		return topic
	}

	topic = findPerOrgnaization(dbConn, appName, org, topicName)

	if topic != nil {
		return topic
	}

	topic = findGlobal(dbConn, topicName)

	if topic != nil {
		return topic
	}

	return nil

}
