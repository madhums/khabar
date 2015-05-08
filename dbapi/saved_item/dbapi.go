package saved_item

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

func Insert(coll string, savedItem *db.SavedItem) string {
	savedItem.PrepareSave()
	return db.Conn.Insert(coll, savedItem)
}

func Get(coll string, query *utils.M) (savedItem *db.SavedItem, err error) {
	savedItem = new(db.SavedItem)

	err = db.Conn.GetOne(coll, *query, savedItem)

	if err != nil {
		return nil, err
	}

	return savedItem, nil
}
