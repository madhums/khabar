/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run 001_modify_available_topics.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Adds `channels` array to `topics_available` collection
 * - Adds `value` property to `topics` collection which determines if the
 * 	 setting is a `default` one on an org level
 */

package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	availableCollection = "topics_available"
	topicsCollection    = "topics"
)

/**
 * Main
 */

func main() {
	session, db := Connect()
	MigrateAvailable(db)
	MigrateTopics(db)
	session.Close()
	fmt.Println("\n", "Closing mongodb connection")
}

/**
 * Connect to mongo
 */

func Connect() (*mgo.Session, *mgo.Database) {
	uri := os.Getenv("MONGODB_URL")

	if uri == "" {
		uri = "mongodb://localhost:27017/notifications_testing"
	}

	mInfo, err := mgo.ParseURL(uri)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	session.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri, "\n")

	sess := session.Clone()

	return session, sess.DB(mInfo.Database)
}

/**
 * Migrate Available
 */

func MigrateAvailable(db *mgo.Database) (err error) {
	Available := db.C(availableCollection)

	change, err := Available.UpdateAll(
		bson.M{},
		bson.M{
			"$set": bson.M{
				"channels": []string{"email", "web", "push"},
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", availableCollection, "` collection")
	return
}

/**
 * Migrate Topics
 */

func MigrateTopics(db *mgo.Database) (err error) {
	Topics := db.C(topicsCollection)

	change, err := Topics.UpdateAll(
		bson.M{
			"user": "",
			"org": bson.M{
				"$ne": "",
			},
		},
		bson.M{
			"$set": bson.M{
				"value": true,
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", topicsCollection, "` collection")
	return
}

/**
 * Handle errors
 */

func handle_errors(err error) {
	if err != nil {
		log.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
