package pending

import (
	"time"

	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

const LATENCY = 2 * int64(time.Minute/time.Nanosecond)

func Throttled(pending_item *PendingItem) bool {
	var query utils.M = make(utils.M)

	query["org"] = pending_item.Organization
	query["app_name"] = pending_item.AppName
	query["topic"] = pending_item.Topic
	query["user"] = pending_item.User
	query["created_by"] = pending_item.CreatedBy
	query["created_on"] = utils.M{"$gt": utils.EpochNow() - LATENCY}

	count := db.Conn.Count(db.SentCollection, query)

	return count > 0
}
