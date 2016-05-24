// +build ignore

/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run 003_clear_org_user_settings.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Removes org and user settings from topics collection
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * Main
 */

func main() {
	session, db, _ := Connect()

	for collection, keys := range db.GetIndexes() {
		index := mgo.Index{
			Key:        strings.Split(keys, " "),
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}

		err := db.C(collection).EnsureIndex(index)
		if err != nil {
			log.Println("Error creating index:", err)
			panic(err)
		}
	}

	Topics := db.C(topicsCollection)
	query := bson.M{
		"org":  utils.M{"$ne": ""},
		"user": utils.M{"$ne": ""},
	}

	// Remove user preference
	change, err := Topics.RemoveAll(query)
	handle_errors(err)
	fmt.Println(change.Removed, "user preferences  were removed from `", topicsCollection, "` collection")

	// Remove org preference
	query["user"] = ""
	change, err = Topics.RemoveAll(query)
	handle_errors(err)
	fmt.Println(change.Removed, "org preferences  were removed from `", topicsCollection, "` collection")

	session.Close()
	fmt.Println("\n", "Closing mongodb connection")
}

/**
 * Connect to mongo
 */

func Connect() (*mgo.Session, *mgo.Database, string) {
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

	return session, sess.DB(mInfo.Database), mInfo.Database
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
