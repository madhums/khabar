// +build ignore

/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run 002_add_defaults_locks_to_topics.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Remove `app_name` property from `topics` collection
 * - Remove previously added `value` property to `topics` collection
 * - Modify `channels` property in `topics` collection from array of strings
 *   to array of objects
 * - Rename locks collection to `temp_locks`
 * - Rename defaults collection to `temp_defaults`
 */

package migrations

import (
	"fmt"
	"log"
	"os"

	models "github.com/bulletind/khabar/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	availableCollection = "topics_available"
	topicsCollection    = "topics"
)

type Topic struct {
	models.BaseModel `bson:",inline"`
	User             string   `json:"user" bson:"user"`
	Organization     string   `json:"org" bson:"org"`
	Channels         []string `json:"channels" bson:"channels" required:"true"`
	Ident            string   `json:"ident" bson:"ident" required:"true"`
}

/**
 * Main
 */

func main() {
	session, db, dbName := Connect()
	RemoveAppName(db)
	RemoveValue(db)
	ModifyChannels(db)
	RemoveLocks(session, dbName)
	RemoveDefaults(session, dbName)
	session.Close()
	fmt.Println("\n", "Closing mongodb connection")
}

/**
 * Remove `app_name` from `topics` collection
 */

func RemoveAppName(db *mgo.Database) (err error) {
	Topics := db.C(topicsCollection)

	change, err := Topics.UpdateAll(
		bson.M{},
		bson.M{
			"$unset": bson.M{
				"app_name": "",
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", topicsCollection, "` collection")
	return
}

/**
 * Remove `value` from `topics` collection
 */

func RemoveValue(db *mgo.Database) (err error) {
	Topics := db.C(topicsCollection)

	change, err := Topics.UpdateAll(
		bson.M{},
		bson.M{
			"$unset": bson.M{
				"value": "",
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", topicsCollection, "` collection")
	return
}

/**
 * Modify channels array in `topics` collection
 */

func ModifyChannels(db *mgo.Database) (err error) {
	Topics := db.C(topicsCollection)
	notNull := bson.M{"$ne": ""}
	query := bson.M{
		"user": notNull,
		"org":  notNull,
	}
	var topics []Topic

	// Modify user setting

	err = Topics.Find(query).All(&topics)
	handle_errors(err)

	for _, topic := range topics {
		channels := make([]models.Channel, 0)

		for _, name := range topic.Channels {
			channel := new(models.Channel)
			channel.Name = name
			channel.Enabled = true
			channels = append(channels, *channel)
		}

		Topics.Update(bson.M{"_id": topic.Id}, bson.M{
			"$set": bson.M{
				"channels": channels,
			},
		})
	}

	// fmt.Println("Updated", change.Updated, "documents in `", topicsCollection, "` collection")

	// Modify org setting

	query["user"] = ""
	change, err := Topics.UpdateAll(
		query,
		bson.M{
			"$set": map[string][]models.Channel{
				"channels": make([]models.Channel, 0),
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", topicsCollection, "` collection")

	return
}

/**
 * Remove locks collection
 */

func RemoveLocks(session *mgo.Session, dbName string) {
	var result interface{}
	err := session.Run(bson.D{
		{"renameCollection", dbName + ".locks"},
		{"to", dbName + ".temp_locks"},
	}, result)
	if err != nil {
		fmt.Println("Error removing locks collection. It may have already been removed")
	}
}

/**
 * Remove defaults collection
 */

func RemoveDefaults(session *mgo.Session, dbName string) {
	var result interface{}
	err := session.Run(bson.D{
		{"renameCollection", dbName + ".defaults"},
		{"to", dbName + ".temp_defaults"},
	}, result)
	if err != nil {
		fmt.Println("Error removing defaults collection. It may have already been removed")
	}
}
