package pending

import (
	"time"

	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
)

//const LATENCY = 1 * int64(time.Second/time.Millisecond)
const LATENCY = 1 * int64(time.Minute/time.Millisecond)

func Throttled(pending_item *db.PendingItem) bool {
	var query utils.M = make(utils.M)

	query["org"] = pending_item.Organization
	query["app_name"] = pending_item.AppName
	query["topic"] = pending_item.Topic
	query["user"] = pending_item.User
	query["created_by"] = pending_item.CreatedBy
	query["created_on"] = utils.M{"$gt": utils.EpochNow() - LATENCY}
	query["entity"] = pending_item.Entity

	count := db.Conn.Count(db.SentCollection, query)

	return count > 0
}
