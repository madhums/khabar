// +build ignore

/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run migrations/005_rename_category.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Rename category 'incidentapp' to 'observationpp'
 * 	 setting is a `default` one on an org level
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bulletind/khabar/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	availableCollection = "topics_available"
)

/**
 * Main
 */

func main() {
	session, db, _ := Connect()
	defer session.Close()
	MigrateAvailable(db, "incidentapp", "observationapp")
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

	const SSL_SUFFIX = "?ssl=true"
	connStringForDb := uri
	if strings.HasSuffix(uri, SSL_SUFFIX) {
		connStringForDb = strings.TrimSuffix(uri, SSL_SUFFIX)
	}

	mInfo, _ := mgo.ParseURL(connStringForDb)
	conn := db.GetConn(uri, mInfo.Database)
	session := conn.Session
	session.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri, "\n")

	sess := session.Clone()

	return session, sess.DB(mInfo.Database), mInfo.Database
}

/**
 * Migrate Available
 */

func MigrateAvailable(db *mgo.Database, fromCategory string, toCategory string) (err error) {
	Available := db.C(availableCollection)

	change, err := Available.UpdateAll(
		bson.M{
			"app_name": fromCategory,
		},
		bson.M{
			"$set": bson.M{
				"app_name": toCategory,
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", availableCollection, "` collection")
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
