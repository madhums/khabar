package saved_item

import "github.com/bulletind/khabar/db"

func Insert(coll string, savedItem *db.SavedItem) string {
	savedItem.PrepareSave()
	return db.Conn.Insert(coll, savedItem)
}
