package stats

import (
	"log"
	"sync"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

type Stats struct {
	LastSeen    int64 `json:"last_seen"`
	TotalCount  int   `json:"total_count"`
	UnreadCount int   `json:"unread_count"`
	TotalUnread int   `json:"total_unread"`
}

type RequestArgs struct {
	Organization string `json:"org"`
	AppName      string `json:"app_name"`
	User         string `json:"user" required:"true"`
}

func Save(args *RequestArgs) error {
	user := args.User
	appName := args.AppName
	org := args.Organization

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

	return db.Conn.Upsert(db.StatsCollection, stats_query, save_doc)
}

func Get(args *RequestArgs) (stats *Stats, err error) {
	user := args.User
	appName := args.AppName
	org := args.Organization

	stats = &Stats{}

	stats_query := utils.M{}
	unread_query := utils.M{"is_read": false}
	unread_since_query := utils.M{"is_read": false}

	stats_query["user"] = user
	unread_query["user"] = user
	unread_since_query["user"] = user

	if len(appName) > 0 {
		stats_query["app_name"] = appName
		unread_query["app_name"] = appName
		unread_since_query["app_name"] = appName
	}

	if len(org) > 0 {
		stats_query["org"] = org
		unread_query["org"] = org
		unread_since_query["org"] = org
	}

	var last_seen db.LastSeen

	err = db.Conn.GetOne(db.StatsCollection, stats_query, &last_seen)
	if err != nil {
		err = Save(args)
		if err == nil {
			err = db.Conn.GetOne(db.StatsCollection, stats_query, &last_seen)
		} else {
			log.Println(err)
			return nil, err
		}
	}

	if last_seen.Timestamp > 0 {
		unread_since_query["created_on"] = utils.M{"$gt": last_seen.Timestamp}
	}

	stats.LastSeen = last_seen.Timestamp

	var wg sync.WaitGroup
	wg.Add(3)

	var unread_count, total_count, total_unread int

	go func() {
		defer wg.Done()
		total_count = db.Conn.Count(db.SentCollection, stats_query)
	}()

	go func() {
		defer wg.Done()
		unread_count = db.Conn.Count(db.SentCollection, unread_since_query)
	}()

	go func() {
		defer wg.Done()
		total_unread = db.Conn.Count(db.SentCollection, unread_query)
	}()

	wg.Wait()

	stats.TotalCount = total_count
	stats.TotalUnread = total_unread
	stats.UnreadCount = unread_count

	return
}
