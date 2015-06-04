package sent

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v3"
)

func Update(id bson.ObjectId, doc *utils.M) error {
	return db.Conn.Update(db.SentCollection, utils.M{
		"_id": id,
	}, utils.M{
		"$set": *doc,
	})
}

func MarkRead(user, appName, org string) error {
	var query utils.M = make(utils.M)

	query["user"] = user

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	doc := utils.M{"$set": utils.M{"is_read": true}}

	return db.Conn.Update(db.SentCollection, query, doc)
}

func GetAll(paginator *gottp.Paginator, user, appName, org string) (*[]db.SentItem, error) {
	var query utils.M = make(utils.M)
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []db.SentItem
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

	session := db.Conn.Session.Copy()
	defer session.Close()

	err := db.Conn.GetCursor(session, db.SentCollection,
		query).Sort("-created_on").Skip(skip).Limit(limit).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Insert(ntfInst *db.SentItem) string {
	return db.Conn.Insert(db.SentCollection, ntfInst)
}
