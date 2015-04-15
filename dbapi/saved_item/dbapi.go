package saved_item

import "github.com/bulletind/khabar/db"

func Insert(channel string, savedItem *db.SavedItem) string {
	savedItem.PrepareSave()
	return db.Conn.Insert("saved_"+channel, savedItem)
}
