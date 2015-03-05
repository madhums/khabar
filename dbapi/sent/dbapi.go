package sent

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v2"
	"log"
)

func Update(dbConn *db.MConn, id bson.ObjectId, doc *db.M) error {
	return dbConn.Update(SentCollection, db.M{
		"_id": id,
	}, db.M{
		"$set": *doc,
	})
}

func MarkRead(dbConn *db.MConn, user string,
	appName string, org string) error {
	var query db.M = make(db.M)

	query["user"] = user

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	log.Println(query)

	doc := db.M{"$set": db.M{"is_read": true}}

	return dbConn.Update(SentCollection, query, doc)

}

func GetAll(dbConn *db.MConn, paginator *gottp.Paginator, user string, appName string, org string) *[]SentItem {
	var query db.M = make(db.M)
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

	dbConn.GetCursor(SentCollection, query).Skip(skip).Limit(limit).All(&result)

	return &result
}

func Insert(dbConn *db.MConn, ntfInst *SentItem) string {
	return dbConn.Insert(SentCollection, ntfInst)
}
