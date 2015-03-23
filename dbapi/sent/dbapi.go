package sent

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v2"
)

func Update(dbConn *db.MConn, id bson.ObjectId, doc *utils.M) error {
	return dbConn.Update(SentCollection, utils.M{
		"_id": id,
	}, utils.M{
		"$set": *doc,
	})
}

func MarkRead(dbConn *db.MConn, user string,
	appName string, org string) error {
	var query utils.M = make(utils.M)

	query["user"] = user

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	doc := utils.M{"$set": utils.M{"is_read": true}}

	return dbConn.Update(SentCollection, query, doc)
}

func GetAll(dbConn *db.MConn, paginator *gottp.Paginator, user string, appName string, org string) (*[]SentItem, error) {
	var query utils.M = make(utils.M)
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []SentItem
	query["user"] = user

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
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

	session := dbConn.Session.Copy()
	defer session.Close()

	err := dbConn.GetCursor(session, SentCollection, query).Skip(skip).Limit(limit).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Insert(dbConn *db.MConn, ntfInst *SentItem) string {
	return dbConn.Insert(SentCollection, ntfInst)
}
