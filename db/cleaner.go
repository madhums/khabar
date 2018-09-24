package db

import (
	"time"

	"github.com/bulletind/khabar/utils"

	"gopkg.in/mgo.v2/bson"
)

func CleanupCollections() {
	tables := []string{ProcessedCollection, SavedEmailCollection, SavedPushCollection, SentCollection}
	rec := new(SentItem)
	dayToKeep := utils.EpochDate(time.Now().AddDate(0, 0, -30))

	go func() {
		session := Conn.Session.Copy()
		defer session.Close()
		for _, table := range tables {
			ids := []bson.ObjectId{}
			cursor := Conn.GetCursor(session, table, utils.M{
				"created_on": utils.M{
					"$lt": dayToKeep,
				},
			}).Select(bson.M{
				"_id": 1, "destination_uri": 1,
			}).Limit(500).Iter()

			for cursor.Next(rec) {
				ids = append(ids, rec.Id)
			}

			Conn.Delete(table, utils.M{
				"_id": utils.M{
					"$in": ids,
				},
			})
		}
	}()
}
