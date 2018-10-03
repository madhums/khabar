// +build ignore

/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run migrations/006_rename_incident_log.go
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
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	topicsCollection    = "topics"
	availableCollection = "topics_available"
)

/**
 * Main
 */

func main() {
	session, db, _ := Connect()
	defer session.Close()

	Migrate(db, "log_used", "observation_used")
	Migrate(db, "log_archived", "observation_archived")
	Migrate(db, "incident_action_completed", "action_completed_casefile")
	Migrate(db, "incident_assigned", "casefile_assigned")
	Migrate(db, "incident_contribution", "casefile_contribution")
	Migrate(db, "log_incoming", "observation_incoming")
	Migrate(db, "high_prio_log_incoming", "observation_incoming_high_prio")
	Migrate(db, "inspection_log_incoming", "observation_incoming_inspection")
	Migrate(db, "inspection_high_prio_log_incoming", "observation_incoming_high_prio_inspection")
	Migrate(db, "high_prio_observation_incoming", "observation_incoming_high_prio")
	Migrate(db, "inspection_observation_incoming", "observation_incoming_inspection")
	Migrate(db, "inspection_high_prio_observation_incoming",
		"observation_incoming_high_prio_inspection")

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

	const SSL_STRING = "ssl=true"
	useSsl := false

	if strings.Contains(uri, SSL_STRING) {
		if strings.Contains(uri, SSL_STRING+"&") {
			uri = strings.Replace(uri, SSL_STRING+"&", "", 1)
		} else {
			uri = strings.Replace(uri, SSL_STRING, "", 1)
		}
		useSsl = true
	}

	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	dialInfo.Timeout = 10 * time.Second

	if useSsl {
		config := tls.Config{}
		config.InsecureSkipVerify = true

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &config)
		}
	}

	// get a mgo session
	session, err := mgo.DialWithInfo(dialInfo)
	sess := session.Clone()

	return session, sess.DB(dialInfo.Database), dialInfo.Database
}

func Migrate(db *mgo.Database, from string, to string) (err error) {
	err = MigrateColl(db, "topics", "ident", from, to)
	if err != nil {
		return err
	}
	err = MigrateColl(db, "topics_available", "ident", from, to)
	return
}

/**
 * Migrate Collection
 */

func MigrateColl(db *mgo.Database, collection string, fieldName string, from string, to string) (err error) {
	change, err := db.C(collection).UpdateAll(
		bson.M{
			fieldName: from,
		},
		bson.M{
			"$set": bson.M{
				fieldName: to,
			},
		},
	)
	handle_errors(err)
	fmt.Println("Updated", change.Updated, "documents in `", collection, "` collection")
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
