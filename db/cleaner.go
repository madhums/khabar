package db

import (
	"regexp"
	"strings"
	"time"

	"github.com/bulletind/khabar/utils"

	"gopkg.in/mgo.v2/bson"
)

func CleanupCollections() {
	session := Conn.Session.Copy()
	defer session.Close()

	tables := []string{ProcessedCollection, SavedEmailCollection, SavedPushCollection, SentCollection}
	rec := new(BaseModel)
	sent := new(SentItem)
	moireHost := utils.GetEnv("MOIRE_HOST", true)
	dayToKeep := utils.EpochDate(time.Now().AddDate(0, 0, -30))

	go func() {
		for _, table := range tables {
			ids := []bson.ObjectId{}
			cursor := Conn.GetCursor(session, table, utils.M{
				"updated_on": utils.M{
					"$lt": dayToKeep,
				},
				"limit": 500,
			}).Select(bson.M{
				"_id": 1,
			}).Sort("updated_on").Iter()

			if table == SentCollection {
				for cursor.Next(sent) {
					ids = append(ids, sent.Id)
					if strings.Contains(sent.DestinationUri, moireHost) {
						re, _ := regexp.Compile("^(" + moireHost + "/assets/)(.*)([?].*)")
						utils.DeleteFile(re.FindStringSubmatch(sent.DestinationUri)[1])
					}
				}
			} else {
				for cursor.Next(rec) {
					ids = append(ids, rec.Id)
				}
			}

			Conn.Delete(table, utils.M{
				"_id": utils.M{
					"$in": ids,
				},
			})
		}
	}()
}
