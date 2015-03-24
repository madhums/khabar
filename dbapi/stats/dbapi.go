package stats

import (
	"log"

	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

type Stats struct {
	LastSeen    int64 `json:"last_seen"`
	TotalCount  int   `json:"total_count"`
	UnreadCount int   `json:"unread_count"`
	TotalUnread int   `json:"total_unread"`
}

func Save(user string, appName string, org string) error {
	db := db.Conn

	stats_query := utils.M{
		"user":     user,
		"app_name": appName,
		"org":      org,
	}

	save_doc := utils.M{
		"user":      user,
		"app_name":  appName,
		"org":       org,
		"timestamp": utils.EpochNow(),
	}

	return dbConn.Upsert(db.StatsCollection, stats_query, save_doc)
}

func Get(user string, appName string, org string) (stats *Stats, err error) {
	db := db.Conn
	stats = &Stats{}

	stats_query := utils.M{}
	unread_query := utils.M{"is_read": false}
	unread_since_query := utils.M{"is_read": false}

	stats_query["user"] = user
	stats_query["app_name"] = appName
	stats_query["org"] = org

	unread_query["user"] = user
	unread_since_query["user"] = user

	if len(appName) > 0 {
		unread_query["app_name"] = appName
		unread_since_query["app_name"] = appName
	}

	if len(org) > 0 {
		unread_query["org"] = org
		unread_since_query["org"] = org
	}

	var last_seen db.LastSeen

	err = dbConn.GetOne(db.StatsCollection, stats_query, &last_seen)
	if err != nil {
		err = Save(dbConn, user, appName, org)
		if err == nil {
			err = dbConn.GetOne(db.StatsCollection, stats_query, &last_seen)
		} else {
			log.Println(err)
			return
		}
	}

	if last_seen.Timestamp > 0 {
		unread_since_query["created_on"] = utils.M{"$gt": last_seen.Timestamp}
	}

	stats.LastSeen = last_seen.Timestamp

	stats.TotalCount = dbConn.Count(db.SentCollection, stats_query)
	stats.UnreadCount = dbConn.Count(db.SentCollection, unread_since_query)
	stats.TotalUnread = dbConn.Count(db.SentCollection, unread_query)

	return
}
